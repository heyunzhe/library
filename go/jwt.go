package mode

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-jwt-secret-key-change-in-production"
	}
	return []byte(secret)
}

var jwtSecret = getJWTSecret()

// Claims 自定义 JWT 声明
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"` // "user" or "admin"
	jwt.RegisteredClaims
}

// GenerateAccessToken 生成访问令牌（15分钟有效期）
func GenerateAccessToken(username, role string) (string, error) {
	claims := Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken() string {
	return uuid.New().String()
}

// ValidateAccessToken 验证访问令牌
func ValidateAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}

// extractToken 从请求中提取令牌（优先 Authorization 头，其次 cookie）
func extractToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// 上下文键类型
type contextKey string

const (
	ContextUsername contextKey = "username"
	ContextRole     contextKey = "role"
)

// JWTAuthMiddleware 验证JWT，适用于普通用户接口
func JWTAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := extractToken(r)
		if tokenStr == "" {
			http.Error(w, "未授权", http.StatusUnauthorized)
			return
		}
		claims, err := ValidateAccessToken(tokenStr)
		if err != nil {
			http.Error(w, "令牌无效或已过期", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ContextUsername, claims.Username)
		ctx = context.WithValue(ctx, ContextRole, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// JWTAdminMiddleware 验证JWT并检查管理员权限
func JWTAdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := extractToken(r)
		if tokenStr == "" {
			http.Error(w, "未授权", http.StatusUnauthorized)
			return
		}
		claims, err := ValidateAccessToken(tokenStr)
		if err != nil {
			http.Error(w, "令牌无效或已过期", http.StatusUnauthorized)
			return
		}
		if claims.Role != "admin" {
			http.Error(w, "权限不足", http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), ContextUsername, claims.Username)
		ctx = context.WithValue(ctx, ContextRole, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GetUsername 从请求上下文中获取用户名
func GetUsername(r *http.Request) string {
	if v := r.Context().Value(ContextUsername); v != nil {
		return v.(string)
	}
	return ""
}

// StoreRefreshToken 存储刷新令牌到数据库
func StoreRefreshToken(username, refreshToken string) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	db.Exec("INSERT INTO refresh_tokens (username, token_hash, expires_at) VALUES (?, ?, ?)",
		username, string(hashed), expiresAt.Format("2006-01-02 15:04:05"))
}

// ValidateRefreshToken 验证刷新令牌
func ValidateRefreshToken(username, refreshToken string) bool {
	var id int
	var tokenHash string
	var expiresAtStr string
	err := db.QueryRow("SELECT id, token_hash, expires_at FROM refresh_tokens WHERE username = ? ORDER BY created_at DESC LIMIT 1",
		username).Scan(&id, &tokenHash, &expiresAtStr)
	if err != nil {
		return false
	}
	expiresAt, err := time.Parse("2006-01-02 15:04:05", expiresAtStr)
	if err != nil {
		return false
	}
	if time.Now().After(expiresAt) {
		db.Exec("DELETE FROM refresh_tokens WHERE id = ?", id)
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(tokenHash), []byte(refreshToken))
	if err != nil {
		return false
	}
	// 使用后删除（一次性刷新令牌）
	db.Exec("DELETE FROM refresh_tokens WHERE id = ?", id)
	return true
}

// RefreshTokenHandler 处理令牌刷新请求
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	refreshToken := r.FormValue("refresh_token")
	if username == "" || refreshToken == "" {
		http.Error(w, "参数不完整", http.StatusBadRequest)
		return
	}

	if !ValidateRefreshToken(username, refreshToken) {
		http.Error(w, "刷新令牌无效或已过期", http.StatusUnauthorized)
		return
	}

	role := "user"
	var adminRole string
	err := db.QueryRow("SELECT admin_role FROM admin WHERE admin_id = ?", username).Scan(&adminRole)
	if err == nil {
		role = "admin"
	}

	newAccessToken, err := GenerateAccessToken(username, role)
	if err != nil {
		http.Error(w, "生成令牌失败", http.StatusInternalServerError)
		return
	}

	newRefreshToken := GenerateRefreshToken()
	StoreRefreshToken(username, newRefreshToken)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"access_token":"` + newAccessToken + `","refresh_token":"` + newRefreshToken + `"}`))
}

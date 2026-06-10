package mode

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func getCookieSecret() []byte {
	secret := os.Getenv("COOKIE_SECRET")
	if secret == "" {
		secret = "default-cookie-secret-key-change-in-production"
	}
	return []byte(secret)
}

var store = sessions.NewCookieStore(getCookieSecret()) // 使用环境变量配置的密钥
var mu sync.Mutex

// 汇总表结构体
type Librarysum struct {
	Total_books_amount  int `json:"total_books_amount"`
	Total_lend_amount   int `json:"total_lend_amount"`
	Total_return_amount int `json:"total_return_amount"`
	Total_users_amount  int `json:"total_users_amount"`
	Cur_user_amount     int `json:"cur_user_amount"`
}

type Rakinng struct {
	Title string `json:"title"`
	Isbn  string `json:"isbn"`
	Cover string `json:"cover"`
	Count int    `json:"count"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "html/library.html", nil)

}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "html/about.html", nil)

}

func generateSessionID() string {
	// 实现一个生成唯一会话ID的函数，例如使用 UUID
	return uuid.New().String()
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "html/login.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		mu.Lock()
		defer mu.Unlock()

		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == "" || password == "" {
			http.Error(w, "账号和密码不能为空", http.StatusBadRequest)
			return
		}

		var storedPassword string
		row := db.QueryRow("SELECT password FROM users WHERE username = ?", username)
		err := row.Scan(&storedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
				errorLog.Println("账号或密码错误：", err)
				return
			}
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			errorLog.Println("服务器错误：", err)
			return
		}

		// bcrypt验证密码
		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
		if err != nil {
			http.Error(w, "密码错误", http.StatusUnauthorized)
			return
		}

		// 保持向后兼容：创建会话（供现有中间件使用）
		session, _ := store.New(r, "user-session")
		sessionID := generateSessionID()
		var existingSessionID string
		_ = db.QueryRow("SELECT session_id FROM session_state WHERE session_name = ?", username).Scan(&existingSessionID)
		if existingSessionID == "" {
			db.Exec("INSERT INTO session_state (session_name, session_id) VALUES (?, ?)", username, sessionID)
		} else {
			db.Exec("UPDATE session_state SET session_id = ? WHERE session_name = ?", sessionID, username)
		}
		session.Values["username"] = username
		session.Values["sessionID"] = sessionID
		session.Values["login"] = true
		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		}
		session.Save(r, w)

		// 生成JWT令牌（供新前端使用）
		accessToken, err := GenerateAccessToken(username, "user")
		if err != nil {
			http.Error(w, "生成令牌失败", http.StatusInternalServerError)
			return
		}
		refreshToken := GenerateRefreshToken()
		StoreRefreshToken(username, refreshToken)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
		return
	}
}

func isValidSession(r *http.Request) (bool, error) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		return false, fmt.Errorf("获取会话失败: %v", err)
	}

	username, ok := session.Values["username"].(string)
	if !ok {
		return false, nil
	}

	sessionID, ok := session.Values["sessionID"].(string)
	if !ok {
		return false, nil
	}

	var storedSessionID string
	err = db.QueryRow("SELECT session_id FROM session_state WHERE session_name = ?", username).Scan(&storedSessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("查询数据库失败: %v", err)
	}

	return sessionID == storedSessionID, nil
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valid, err := isValidSession(r)
		if err != nil {
			errorLog.Println("验证会话时出错:", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}
		if !valid {
			session, err := store.Get(r, "user-session")
			if err != nil {
				http.Error(w, "获取会话失败", http.StatusInternalServerError)
				errorLog.Println("获取会话错误", err)
				return
			}

			delete(session.Values, "username")
			delete(session.Values, "login")
			delete(session.Values, "sessionID")

			session.Options.MaxAge = -1
			err = session.Save(r, w)
			if err != nil {
				http.Error(w, "保存会话失败", http.StatusInternalServerError)
				errorLog.Println("保存会话失败", err)
				return
			}

			// 清除 Cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "user-session",
				Value:    "",
				Path:     "/",
				MaxAge:   -1, // 过期
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			})

			http.Redirect(w, r, "/login", http.StatusSeeOther)

			return
		}
		// 从 session 获取 username 并写入 context，与 JWTAuthMiddleware 行为一致
		session, _ := store.Get(r, "user-session")
		if username, ok := session.Values["username"].(string); ok {
			ctx := context.WithValue(r.Context(), ContextUsername, username)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	}
}

func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	// 获取会话
	session, err := store.Get(r, "user-session")
	username, _ := session.Values["username"].(string)
	if err != nil {
		http.Error(w, "获取会话失败", http.StatusInternalServerError)
		errorLog.Println("获取会话错误", err)
		return
	}

	// 清除会话中的用户信息
	delete(session.Values, "username")
	delete(session.Values, "login")
	delete(session.Values, "sessionID")

	_, err = db.Exec("DELETE FROM session_state WHERE session_name = ?", username)
	if err != nil {
		errorLog.Println("删除会话失败：", err)
		return
	}

	// 过期会话（设置 MaxAge 为 -1）
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "保存会话失败", http.StatusInternalServerError)
		errorLog.Println("保存会话失败", err)
		return
	}

	// 清除 Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "user-session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // 过期
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	// 返回成功响应
	w.WriteHeader(http.StatusOK)
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	// GET 请求渲染后台首页
	if r.Method == http.MethodGet {
		valid, err := AdminisValidSession(r)
		if err != nil {
			errorLog.Println("验证会话时出错:", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}
		if !valid {
			http.Redirect(w, r, "/index", http.StatusSeeOther)
			return
		}
		var librarysum Librarysum

		var sum1, sum2, sum3 int
		err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&sum1)
		if err != nil {
			errorLog.Println("数据库错误：", err)
			return
		}
		err = db.QueryRow("SELECT COUNT(*) FROM admin").Scan(&sum2)
		if err != nil {
			errorLog.Println("数据库错误：", err)
			return
		}
		err = db.QueryRow("SELECT COUNT(*) FROM session_state").Scan(&sum3)
		if err != nil {
			errorLog.Println("数据库错误：", err)
			return
		}
		librarysum.Cur_user_amount = sum3

		_, err = db.Exec("UPDATE library_summary SET Total_users_amount = ?", sum1+sum2)
		if err != nil {
			errorLog.Println("数据库错误：", err)
			return
		}

		row := db.QueryRow("SELECT * FROM library_summary")
		err = row.Scan(&librarysum.Total_books_amount, &librarysum.Total_lend_amount, &librarysum.Total_return_amount, &librarysum.Total_users_amount)
		if err != nil {
			errorLog.Println("数据库错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		renderTemplate(w, "html/admin.html", librarysum)
		return
	}

	// POST 请求处理管理员登录
	if r.Method == http.MethodPost {
		mu.Lock()
		defer mu.Unlock()

		username := r.FormValue("adminID")
		password := r.FormValue("adminPassword")
		if username == "" || password == "" {
			http.Error(w, "账号和密码不能为空", http.StatusBadRequest)
			return
		}

		var storedPassword string
		err := db.QueryRow("SELECT admin_password FROM admin WHERE admin_id = ?", username).Scan(&storedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "账号或密码错误", http.StatusUnauthorized)
				return
			}
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			errorLog.Println("查询管理员失败:", err)
			return
		}

		// 管理员密码（假设存储的是明文或简单加密）
		if storedPassword != password {
			http.Error(w, "密码错误", http.StatusUnauthorized)
			return
		}

		// 创建会话
		session, _ := store.New(r, "admin-session")
		sessionID := generateSessionID()
		db.Exec("INSERT INTO session_state (session_name, session_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE session_id = ?", username, sessionID, sessionID)
		session.Values["adminID"] = username
		session.Values["sessionID"] = sessionID
		session.Values["loggedin"] = true
		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		}
		session.Save(r, w)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "登录成功"})
		return
	}
}

func AdminisValidSession(r *http.Request) (bool, error) {
	session, err := store.Get(r, "admin-session")
	if err != nil {
		return false, fmt.Errorf("获取会话失败: %v", err)
	}

	adminID, ok := session.Values["adminID"].(string)
	if !ok {
		return false, nil
	}

	loggedin, ok := session.Values["loggedin"].(bool)
	if !ok || !loggedin {
		return false, nil
	}

	sessionID, ok := session.Values["sessionID"].(string)
	if !ok {
		return false, nil
	}

	var storedSessionID string
	err = db.QueryRow("SELECT session_id FROM session_state WHERE session_name = ?", adminID).Scan(&storedSessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("查询数据库失败: %v", err)
	}

	return sessionID == storedSessionID, nil
}

func AdminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		valid, err := AdminisValidSession(r)
		if err != nil {
			errorLog.Println("验证会话时出错:", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}
		if !valid {
			// 会话无效，清理 cookie 并重定向到首页
			session, _ := store.Get(r, "admin-session")
			delete(session.Values, "adminID")
			delete(session.Values, "loggedin")
			delete(session.Values, "sessionID")
			session.Options.MaxAge = -1
			session.Save(r, w)

			http.SetCookie(w, &http.Cookie{
				Name:     "admin-session",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			})
			http.Redirect(w, r, "/index", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// 获取会话
	session, err := store.Get(r, "admin-session")
	adminID, _ := session.Values["adminID"].(string)
	if err != nil {
		http.Error(w, "获取会话失败", http.StatusInternalServerError)
		errorLog.Println("获取会话错误", err)
		return
	}

	// 清除会话中的用户信息
	delete(session.Values, "adminID")
	delete(session.Values, "loggedin")
	delete(session.Values, "sessionID")

	_, err = db.Exec("DELETE FROM session_state WHERE session_name = ?", adminID)
	if err != nil {
		errorLog.Println("删除会话失败：", err)
		return
	}

	// 过期会话（设置 MaxAge 为 -1）
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "保存会话失败", http.StatusInternalServerError)
		errorLog.Println("保存会话失败", err)
		return
	}
	// 清除 Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "admin-session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // 过期
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/index", http.StatusSeeOther)
}

func RankingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "html/ranking.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		rows, err := db.Query("SELECT b.title,b.isbn,b.cover,count(*) FROM lend_records l JOIN all_books b ON l.isbn = b.isbn GROUP BY b.isbn")
		if err != nil {
			errorLog.Println("数据库错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var rankings []Rakinng
		for rows.Next() {
			var ranking Rakinng
			err = rows.Scan(&ranking.Title, &ranking.Isbn, &ranking.Cover, &ranking.Count)
			if err != nil {
				errorLog.Println("数据库错误：", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			rankings = append(rankings, ranking)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rankings)
	}
}

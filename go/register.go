package mode

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SendVerifyHandler 发送验证码
func SendVerifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "邮箱不能为空"})
		return
	}

	// 检查邮箱是否已被注册
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&exists)
	if err == nil && exists > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "该邮箱已被注册"})
		return
	}

	// 检查发送频率（60秒内不能重复发送）
	var lastTime time.Time
	err = db.QueryRow("SELECT created_at FROM verification_codes WHERE email = ? ORDER BY created_at DESC LIMIT 1", email).Scan(&lastTime)
	if err == nil && time.Since(lastTime).Seconds() < 60 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]string{"error": "请60秒后再试"})
		return
	}
	// 生成验证码
	code := GenerateVerificationCode()
	expiresAt := time.Now().Add(10 * time.Minute)

	// 删除旧的验证码
	_, err = db.Exec("DELETE FROM verification_codes WHERE email = ?", email)
	if err != nil {
		errorLog.Println("删除旧验证码失败:", err)
	}

	_, err = db.Exec("INSERT INTO verification_codes (email, code, expires_at) VALUES (?, ?, ?)",
		email, code, expiresAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		errorLog.Println("存储验证码失败:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "服务器错误"})
		return
	}

	// 发送邮件
	err = SendVerificationCode(email, code)
	if err != nil {
		errorLog.Println("发送验证码失败:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "验证码发送失败"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "验证码已发送到您的邮箱"})
}

// RegisterHandler 用户注册
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		password := r.FormValue("password")
		email := r.FormValue("email")
		code := r.FormValue("code")

		// 基本验证
		if password == "" || email == "" || code == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "所有字段均为必填"})
			return
		}

		if len(password) < 6 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "密码长度不能少于6位"})
			return
		}

		// 用邮箱作为读者号
		username := email

		// 检查用户名是否已存在
		var exists int
		err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&exists)
		if err == nil && exists > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "该邮箱已被注册"})
			return
		}

		// 验证码校验
		var storedCode string
		var expiresAt time.Time
		err = db.QueryRow("SELECT code, expires_at FROM verification_codes WHERE email = ? ORDER BY created_at DESC LIMIT 1",
			email).Scan(&storedCode, &expiresAt)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "验证码未发送或已过期"})
			return
		}

		if time.Now().After(expiresAt) {
			errorLog.Printf("验证码已过期: now=%v, expires=%v", time.Now(), expiresAt)
			db.Exec("DELETE FROM verification_codes WHERE email = ?", email)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "验证码已过期"})
			return
		}

		// 验证码匹配
		if storedCode != code {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "验证码错误"})
			return
		}

		// 验证通过，删除已使用的验证码
		db.Exec("DELETE FROM verification_codes WHERE email = ?", email)

		// 密码加密
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			errorLog.Println("密码加密失败:", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "服务器错误"})
			return
		}

		// 插入用户（姓名默认同邮箱）
		_, err = db.Exec("INSERT INTO users (name, username, password, email, email_verified) VALUES (?, ?, ?, ?, 1)",
			email, username, string(hashedPassword), email)
		if err != nil {
			errorLog.Println("注册用户失败:", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "注册失败"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message":  "注册成功",
			"username": username,
			"email":    email,
		})
		return
	}

	http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
}

// SendResetCodeHandler 发送重置密码验证码
func SendResetCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "邮箱不能为空"})
		return
	}

	// 检查邮箱是否存在
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&exists)
	if err != nil || exists == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "该邮箱未注册"})
		return
	}

	code := GenerateVerificationCode()
	expiresAt := time.Now().Add(10 * time.Minute)

	_, err = db.Exec("DELETE FROM verification_codes WHERE email = ?", email)
	if err != nil {
		errorLog.Println("删除旧验证码失败:", err)
	}

	_, err = db.Exec("INSERT INTO verification_codes (email, code, expires_at) VALUES (?, ?, ?)",
		email, code, expiresAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		errorLog.Println("存储验证码失败:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "服务器错误"})
		return
	}

	subject := "密码重置验证码"
	body := fmt.Sprintf(`您好！

您正在重置图书馆管理系统的密码。

您的验证码是：%s

该验证码有效期为10分钟。

如果这不是您本人的操作，请忽略此邮件。`, code)
	err = SendEmail(email, subject, body)
	if err != nil {
		errorLog.Println("发送验证码失败:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "验证码发送失败"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "验证码已发送到您的邮箱"})
}

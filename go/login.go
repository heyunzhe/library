package mode

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"net/http"
	"sync"
	"text/template"
)

var store = sessions.NewCookieStore([]byte("your-secret-key")) // 使用安全的密钥
var isAdminLoggedIn bool                                       // 用于追踪是否已有管理员登录
var mu sync.Mutex

// 汇总表结构体
type Librarysum struct {
	Total_books_amount  int `json:"total_books_amount"`
	Total_lend_amount   int `json:"total_lend_amount"`
	Total_return_amount int `json:"total_return_amount"`
	Total_users_amount  int `json:"total_users_amount"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("html/library.html")
		if err != nil {
			fmt.Printf("解析模板失败: %v\n", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
		}
		err = tmpl.ExecuteTemplate(w, "library.html", nil)
		if err != nil {
			fmt.Printf("执行模板失败: %v\n", err)
			errorLog.Println("服务器错误：", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
		}
	}
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("html/about.html")
		if err != nil {
			fmt.Printf("解析模板失败: %v\n", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
		}
		err = tmpl.ExecuteTemplate(w, "about.html", nil)
		if err != nil {
			fmt.Printf("执行模板失败: %v\n", err)
			errorLog.Println("服务器错误：", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
		}
	}
}

func generateSessionID() string {
	// 实现一个生成唯一会话ID的函数，例如使用 UUID
	return uuid.New().String()
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		tmpl, err := template.ParseFiles("html/login.html")
		if err != nil {
			fmt.Printf("解析模板失败: %v\n", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
		}
		err = tmpl.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			fmt.Printf("执行模板失败: %v\n", err)
			errorLog.Println("服务器错误：", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
		}
	}

	if r.Method == http.MethodPost {
		mu.Lock()
		defer mu.Unlock()

		username := r.FormValue("username")
		password := r.FormValue("password")
		var user Users
		row := db.QueryRow("SELECT username , password FROM users WHERE username = ?", username)
		err := row.Scan(&user.Username, &user.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized) //401
				errorLog.Println("帐号或密码错误：", err)
				return
			}
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			errorLog.Println("服务器错误：", err)
			return
		}

		if username == user.Username && password == user.Password {
			session, _ := store.New(r, "user-session")
			sessionID := generateSessionID()
			var existingSessionID string
			err = db.QueryRow("SELECT session_id FROM session_state WHERE session_name = ?", username).Scan(&existingSessionID)
			if err == sql.ErrNoRows {
				_, err = db.Exec("INSERT INTO session_state (session_name, session_id) VALUES (?, ?)", username, sessionID)
				if err != nil {
					errorLog.Println("数据库错误：", err)
					return
				}
			} else if err != nil {
				errorLog.Println("查询错误：", err)
				return
			} else {
				_, err = db.Exec("UPDATE session_state SET session_id = ? WHERE session_name = ?", sessionID, username)
				if err != nil {
					errorLog.Println("数据库错误：", err)
					return
				}
			}

			session.Values["username"] = username
			session.Values["sessionID"] = sessionID
			session.Values["login"] = true
			session.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   3600, // 设置 cookie 过期时间
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			}

			err = session.Save(r, w)
			if err != nil {
				http.Error(w, "保存会话失败", http.StatusInternalServerError)
				errorLog.Println("保存会话失败", err)
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		}
		http.Error(w, "密码错误", http.StatusUnauthorized)
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
	mu.Lock()
	defer mu.Unlock()

	session, err := store.Get(r, "admin-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		errorLog.Println("获取会话错误", err)
		return
	}

	if r.Method == http.MethodPost {
		adminPassword := r.FormValue("adminPassword")
		if adminPassword == "123456" {
			if isAdminLoggedIn {
				http.Error(w, "已有其他管理员登录", http.StatusUnauthorized)
				return
			}
			session.Values["admin"] = "admin"
			session.Values["loggedin"] = true
			isAdminLoggedIn = true // 标记为有管理员登录

			session.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   3600, // 设置 cookie 过期时间
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			}

			err = session.Save(r, w)
			if err != nil {
				http.Error(w, "保存会话失败", http.StatusInternalServerError)
				errorLog.Println("保存会话失败", err)
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		}
		http.Error(w, "密码错误", http.StatusUnauthorized)
		return
	}

	loggedin, ok := session.Values["loggedin"].(bool)
	if !ok || !loggedin {
		http.Error(w, "未登录", http.StatusUnauthorized)
		return
	}

	var librarysum Librarysum
	row := db.QueryRow("SELECT * FROM library_summary")
	err = row.Scan(&librarysum.Total_books_amount, &librarysum.Total_lend_amount, &librarysum.Total_return_amount, &librarysum.Total_users_amount)
	if err != nil {
		errorLog.Println("数据库错误：", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("html/admin.html")
	if err != nil {
		fmt.Printf("解析模板失败: %v\n", err)
		http.Error(w, "服务器错误", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, librarysum)
	if err != nil {
		errorLog.Println("编码错误：", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func AdminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "admin-session")
		if err != nil {
			http.Error(w, "会话错误", http.StatusInternalServerError)
			errorLog.Println("获取会话错误", err)
			return
		}

		loggedin, ok := session.Values["loggedin"].(bool)
		if !ok || !loggedin {
			http.Error(w, "权限不足", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	session, err := store.Get(r, "admin-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		errorLog.Println("获取会话错误", err)
		return
	}

	session.Values = make(map[interface{}]interface{})
	// 设置 Cookie 的过期时间为过去的时间
	session.Options.MaxAge = -1
	session.Values["loggedin"] = false
	delete(session.Values, "admin") // 清除管理员信息
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "保存会话失败", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "admin-session", // 确保这里的名称与设置 cookie 时一致
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	isAdminLoggedIn = false // 清除管理员登录状态
	http.Redirect(w, r, "/index", http.StatusSeeOther)
}

func RankingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("html/ranking.html")
		if err != nil {
			fmt.Printf("解析模板失败: %v\n", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
		}
		err = tmpl.ExecuteTemplate(w, "ranking.html", nil)
		if err != nil {
			fmt.Printf("执行模板失败: %v\n", err)
			errorLog.Println("服务器错误：", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
		}
	}
}

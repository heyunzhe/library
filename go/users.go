package mode

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// 用户结构体
type Users struct {
	Name                 string  `json:"name"`
	Username             string  `json:"username"`
	Password             string  `json:"password"`
	User_cur_lend_amount int     `json:"user_cur_lend_amount"`
	User_his_lend_amount int     `json:"user_his_lend_amount"`
	Birthday             string  `json:"birthday"`
	Age                  int     `json:"age"`
	Photo                string  `json:"photo"`
	Email                string  `json:"email"`
	Email_verified       int     `json:"email_verified"`
	Created_at           string  `json:"created_at"`
	Title                string  `json:"title"`
	ISBN                 string  `json:"isbn"`
	Lend_date            string  `json:"lend-date"`
	Exp_return_date      string  `json:"exp_return_date"`
	Return_date          string  `json:"return_date"`
	Late_fee             float64 `json:"late_fee"`
}

type UserInfo struct {
	Name                 string `json:"name"`
	Username             string `json:"username"`
	User_cur_lend_amount int    `json:"user_cur_lend_amount"`
	User_his_lend_amount int    `json:"user_his_lend_amount"`
	Birthday             string `json:"birthday"`
	Age                  int    `json:"age"`
	Photo                string `json:"photo"`
}

type CurrentLoan struct {
	Title           string `json:"title"`
	ISBN            string `json:"isbn"`
	Lend_date       string `json:"lend_date"`
	Exp_return_date string `json:"exp_return_date"`
}

type LoanHistory struct {
	Title       string  `json:"title"`
	ISBN        string  `json:"isbn"`
	Lend_date   string  `json:"lend_date"`
	Return_date string  `json:"return_date"`
	Late_fee    float64 `json:"late_fee"`
}

// 用户意见建议结构体
type Useropi struct {
	Opinion_id int    `json:"opinion_id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Idea       string `json:"idea"`
}

// type Userlend struct {
// 	Lend_id         int    `json:"lend_id"`
// 	Username        string `json:"username"`
// 	Title           string `json:"title"`
// 	ISBN            string `json:"isbn"`
// 	Lend_date       string `json:"lend-date"`
// 	Exp_return_date string `json:"exp_return_date"`
// }

// 查询用户
func ViewUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		session, _ := store.Get(r, "admin-session")
		adminID, _ := session.Values["adminID"].(string)

		var role string
		err := db.QueryRow("SELECT admin_role FROM admin WHERE admin_id = ?", adminID).Scan(&role)
		if err != nil {
			errorLog.Println("数据库错误：", err)
			return
		}
		renderTemplate(w, "html/view-user.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username") //获取输入读者号信息
		user_name := "%" + username + "%"   //用于模糊查询
		//查询读者号符合的用户信息
		row, err := db.Query("SELECT * FROM users WHERE username like ?", user_name)
		if err != nil {
			errorLog.Println("数据库错误", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer row.Close()
		var users []Users
		for row.Next() {
			var user Users
			err = row.Scan(&user.Name, &user.Username, &user.Password, &user.User_cur_lend_amount, &user.User_his_lend_amount, &user.Birthday, &user.Age, &user.Photo, &user.Email, &user.Email_verified, &user.Created_at)
			if err != nil {
				errorLog.Println("服务器错误：", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			errorLog.Println("编码错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
}

// 查看汇总信息
func SumlibraryHandler(w http.ResponseWriter, r *http.Request) {

	var librarysum Librarysum
	row := db.QueryRow("SELECT * FROM library_summary")
	err := row.Scan(&librarysum.Total_books_amount, &librarysum.Total_lend_amount, &librarysum.Total_return_amount, &librarysum.Total_users_amount)
	if err != nil {
		errorLog.Println("数据库错误：", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("html/admin.html")
	if err != nil {
		errorLog.Println("模板解析错误：", err)
		http.Error(w, "内部服务器错误", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, librarysum)
	// w.Header().Set("Content-Type", "application/json")
	// err = json.NewEncoder(w).Encode(librarysum)
	if err != nil {
		errorLog.Println("编码错误：", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// 查看用户意见建议
func ViewUserOpinionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		opinionid := r.FormValue("opinion_id")
		opinion_id := "%" + opinionid + "%"

		rows, err := db.Query("SELECT * FROM user_opinions WHERE opinion_id LIKE ?", opinion_id)
		if err != nil {
			errorLog.Println("服务器错误", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var useropis []Useropi
		for rows.Next() {
			var useropi Useropi
			err := rows.Scan(&useropi.Opinion_id, &useropi.Name, &useropi.Phone, &useropi.Email, &useropi.Idea)
			if err != nil {
				errorLog.Println("服务器错误：", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			useropis = append(useropis, useropi) //循环累加到切片
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(useropis)
		if err != nil {
			errorLog.Println("编码错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}

	if r.Method == http.MethodGet {
		session, _ := store.Get(r, "admin-session")
		adminID, _ := session.Values["adminID"].(string)

		var role string
		err := db.QueryRow("SELECT admin_role FROM admin WHERE admin_id = ?", adminID).Scan(&role)
		if err != nil {
			errorLog.Println("数据库错误：", err)
			return
		}
		renderTemplate(w, "html/view-useropi.html", nil)
		}
}

// 回复用户意见建议
func ReplayUserOpinionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		mu.Lock()
		defer mu.Unlock()
		replay_name := r.FormValue("replay_name")
		replay_date := r.FormValue("replay_date")
		replay_idea := r.FormValue("replay_idea")
		replay_user := r.FormValue("replay_user")

		nowreplaydate, err := time.Parse("2006-01-02", replay_date) // 确保格式与数据库中的一致
		if err != nil {
			errorLog.Println("解析日期错误:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		now := time.Now().Truncate(24 * time.Hour)

		day := nowreplaydate.Sub(now.Truncate(24*time.Hour)).Hours() / 24

		if day == 0 {
			_, err = db.Exec("INSERT INTO replay_opinions (replay_name,replay_date,replay_idea,replay_user) values(?,?,?,?)", replay_name, replay_date, replay_idea, replay_user)
			if err != nil {
				errorLog.Println("数据库错误：", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

// 用户上传意见
func AddUserOpinionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		mu.Lock()
		name := r.FormValue("name")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		idea := r.FormValue("idea")

		_, err := db.Exec("INSERT INTO user_opinions (name,phone,email,idea) values(?,?,?,?)", name, phone, email, idea)

		if err != nil {
			errorLog.Println("服务器错误", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		mu.Unlock()
	}
}

// 个人中心
func UserLibraryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		session, _ := store.Get(r, "user-session")
		_, ok := session.Values["username"].(string)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		renderTemplate(w, "html/user-library.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		session, _ := store.Get(r, "user-session")
		user, ok := session.Values["username"].(string)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var response struct {
			UserInfo     UserInfo      `json:"user_info"`
			CurrentLoans []CurrentLoan `json:"current_loans"`
			LoanHistory  []LoanHistory `json:"loan_history"`
		}

		// 获取用户信息
		err := db.QueryRow("SELECT name, username, user_cur_lend_amount, user_his_lend_amount, birthday, age, photo FROM users WHERE username = ?", user).Scan(
			&response.UserInfo.Name,
			&response.UserInfo.Username,
			&response.UserInfo.User_cur_lend_amount,
			&response.UserInfo.User_his_lend_amount,
			&response.UserInfo.Birthday,
			&response.UserInfo.Age,
			&response.UserInfo.Photo,
		)
		if err != nil {
			errorLog.Println("数据库错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// 获取当前借阅记录
		rows, err := db.Query("SELECT title, isbn, lend_date, exp_return_date FROM cur_lend_records WHERE username = ?", user)
		if err != nil {
			errorLog.Println("数据库错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var loan CurrentLoan
			err = rows.Scan(&loan.Title, &loan.ISBN, &loan.Lend_date, &loan.Exp_return_date)
			if err != nil {
				errorLog.Println("服务器错误：", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			response.CurrentLoans = append(response.CurrentLoans, loan)
		}

		// 获取借阅历史
		hisrows, err := db.Query("SELECT title, isbn, lend_date, return_date, late_fee FROM return_records WHERE username = ?", user)
		if err != nil {
			errorLog.Println("数据库错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer hisrows.Close()
		for hisrows.Next() {
			var historyLoan LoanHistory
			err = hisrows.Scan(&historyLoan.Title, &historyLoan.ISBN, &historyLoan.Lend_date, &historyLoan.Return_date, &historyLoan.Late_fee)
			if err != nil {
				errorLog.Println("服务器错误：", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			response.LoanHistory = append(response.LoanHistory, historyLoan)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// 更新用户信息
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "user-session")
	user, ok := session.Values["username"].(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	name := r.FormValue("name")
	password := r.FormValue("password")
	newpassword := r.FormValue("newpassword")
	birthday := r.FormValue("birthday")
	avatarPath := r.FormValue("avatar_path")

	age := 0
	if birthday != "" {
		b, err := time.Parse("2006-01-02", birthday)
		if err == nil {
			age = time.Now().Year() - b.Year()
			if time.Now().Month() < b.Month() || (time.Now().Month() == b.Month() && time.Now().Day() < b.Day()) {
				age--
			}
			if age < 0 {
				age = 0
			}
		}
	}

	var photoPath string
	file, fileHeader, err := r.FormFile("photo")
	if err == nil {
		defer file.Close()
		ext := ".jpg"
		if fileHeader != nil && fileHeader.Filename != "" {
			if strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".png") {
				ext = ".png"
			} else if strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".gif") {
				ext = ".gif"
			}
		}
		filename := "image-" + user + ext
		photoPath = filepath.Join("images", filename)

		outfile, cerr := os.Create(photoPath)
		if cerr != nil {
			errorLog.Println("无法创建文件:", cerr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer outfile.Close()
		io.Copy(outfile, file)
	} else if avatarPath != "" {
		photoPath = strings.TrimPrefix(avatarPath, "/")
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "服务器错误", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	if password != "" {
		var storedPwd string
		err = tx.QueryRow("SELECT password FROM users WHERE username = ?", user).Scan(&storedPwd)
		if err != nil {
			errorLog.Println("数据库错误:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(storedPwd), []byte(password)) != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "当前密码错误"})
			return
		}
		if newpassword != "" {
			hashed, _ := bcrypt.GenerateFromPassword([]byte(newpassword), bcrypt.DefaultCost)
			if photoPath != "" {
				_, err = tx.Exec("UPDATE users SET name=?, password=?, age=?, birthday=?, photo=? WHERE username=?", name, string(hashed), age, birthday, photoPath, user)
			} else {
				_, err = tx.Exec("UPDATE users SET name=?, password=?, age=?, birthday=? WHERE username=?", name, string(hashed), age, birthday, user)
			}
		} else {
			if photoPath != "" {
				_, err = tx.Exec("UPDATE users SET name=?, age=?, birthday=?, photo=? WHERE username=?", name, age, birthday, photoPath, user)
			} else {
				_, err = tx.Exec("UPDATE users SET name=?, age=?, birthday=? WHERE username=?", name, age, birthday, user)
			}
		}
	} else {
		if photoPath != "" {
			_, err = tx.Exec("UPDATE users SET name=?, age=?, birthday=?, photo=? WHERE username=?", name, age, birthday, photoPath, user)
		} else {
			_, err = tx.Exec("UPDATE users SET name=?, age=?, birthday=? WHERE username=?", name, age, birthday, user)
		}
	}
	if err != nil {
		errorLog.Println("数据库错误:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		errorLog.Println("提交事务失败:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "保存成功"})
}

// generateRandomPassword 生成随机密码
func generateRandomPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 12)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}

// 重置用户密码（通过验证码）
func ResetpasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email") // 前端传递邮箱
		if email == "" {
			http.Error(w, "邮箱不能为空", http.StatusBadRequest)
			return
		}

		// 获取最新的验证码
		var storedCode string
		var expiresAt time.Time
		err := db.QueryRow("SELECT code, expires_at FROM verification_codes WHERE email = ? ORDER BY created_at DESC LIMIT 1", email).Scan(&storedCode, &expiresAt)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "验证码未发送或已过期"})
			return
		}

		if time.Now().After(expiresAt) {
			db.Exec("DELETE FROM verification_codes WHERE email = ?", email)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "验证码已过期"})
			return
		}

		// 验证码校验
		code := r.FormValue("code")
		if storedCode != code {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "验证码错误"})
			return
		}

		// 验证通过，删除已使用的验证码
		db.Exec("DELETE FROM verification_codes WHERE email = ?", email)

		// 检查邮箱是否存在
		var exists int
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&exists)
		if err != nil || exists == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "该邮箱未注册"})
			return
		}

		// 使用前端传递的新密码（通过code字段已废弃，现在使用password）
		password := r.FormValue("password")
		if password == "" {
			// 如果没有传递新密码，生成随机密码
			password = generateRandomPassword()
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			errorLog.Println("密码加密失败:", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("UPDATE users SET password = ? WHERE email = ? ", string(hashed), email)
		if err != nil {
			errorLog.Println("数据库错误:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":  "密码重置成功",
			"password": password, // 返回密码（生产环境移除）
		})
	}
}

// AdminResetPasswordHandler 管理员重置密码（无需验证码）
func AdminResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		if username == "" {
			http.Error(w, "用户名不能为空", http.StatusBadRequest)
			return
		}

		// 生成随机密码
		newPassword := generateRandomPassword()
		hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			errorLog.Println("密码加密失败:", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("UPDATE users SET password = ? WHERE username = ? ", string(hashed), username)
		if err != nil {
			errorLog.Println("数据库错误:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "密码已重置",
			"password": newPassword,
		})
	}
}

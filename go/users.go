package mode

import (
	// "GoPath/librarys/library"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
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
		if role == "Admin" {
			tmpl, err := template.ParseFiles("html/view-user.html")
			if err != nil {
				fmt.Printf("解析模板失败: %v\n", err)
				http.Error(w, "服务器错误", http.StatusInternalServerError)
			}
			err = tmpl.ExecuteTemplate(w, "view-user.html", nil)
			if err != nil {
				fmt.Printf("执行模板失败: %v\n", err)
				errorLog.Println("服务器错误：", err)
				http.Error(w, "服务器错误", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "权限不足", http.StatusUnauthorized)
			return
		}
	}
	if r.Method == http.MethodPost {
		username := r.FormValue("username") //获取输入读者号信息
		user_name := "%" + username + "%"   //用于模糊查询
		// if username != "" {
		//查询读者号符合的用户信息
		row, err := db.Query("SELECT * FROM users WHERE username like ?", user_name)
		if err == sql.ErrNoRows { //判断是否为空
			errorLog.Println("没有这个用户")
			w.WriteHeader(http.StatusNotFound)
			if err != nil { //其他错误
				errorLog.Println("数据库错误", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		defer row.Close()
		var users []Users
		for row.Next() {
			var user Users
			err = row.Scan(&user.Name, &user.Username, &user.Password, &user.User_cur_lend_amount, &user.User_his_lend_amount, &user.Birthday, &user.Age, &user.Photo)
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
		if role == "Admin" {
			tmpl, err := template.ParseFiles("html/view-useropi.html")
			if err != nil {
				fmt.Printf("解析模板失败: %v\n", err)
				http.Error(w, "服务器错误", http.StatusInternalServerError)
			}
			err = tmpl.ExecuteTemplate(w, "view-useropi.html", nil)
			if err != nil {
				fmt.Printf("执行模板失败: %v\n", err)
				errorLog.Println("服务器错误：", err)
				http.Error(w, "服务器错误", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "权限不足", http.StatusUnauthorized)
			return
		}
	}
}

// 回复用户意见建议
func ReplayUserOpinionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		mu.Lock()
		replay_name := r.FormValue("replay_name")
		replay_date := r.FormValue("replay_date")
		replay_idea := r.FormValue("replay_idea")
		replay_user := r.FormValue("replay_user")

		_, err := db.Exec("INSERT INTO replay_opinions (replay_name,replay_date,replay_idea,replay_user) values(?,?,?,?)", replay_name, replay_date, replay_idea, replay_user)
		if err != nil {
			errorLog.Println("数据库错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		mu.Unlock()
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
	session, _ := store.Get(r, "user-session")
	user, ok := session.Values["username"].(string)

	// user := r.FormValue("user")

	if r.Method == http.MethodGet {
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tmpl, err := template.ParseFiles("html/user-library.html")
		if err != nil {
			fmt.Printf("解析模板失败: %v\n", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}
		err = tmpl.ExecuteTemplate(w, "user-library.html", nil)
		if err != nil {
			fmt.Printf("执行模板失败: %v\n", err)
			errorLog.Println("服务器错误：", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
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
	session, _ := store.Get(r, "user-session")
	user, ok := session.Values["username"].(string)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	mu.Lock()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "服务器错误", http.StatusInternalServerError)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback() // 在出现错误时回滚
		}
	}()
	name := r.FormValue("name")
	password := r.FormValue("password")
	newpassword := r.FormValue("newpassword")
	age := r.FormValue("age")
	birthday := r.FormValue("birthday")

	file, _, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "无法获取文件", http.StatusBadRequest)
		errorLog.Println("无法获取文件", err)
		return
	}
	defer file.Close()
	filename := "image-" + user + ".jpg"          //文件命名
	filepath := filepath.Join("images", filename) //保存图片文件

	fileChan := make(chan error)

	go func() {
		outfile, err := os.Create(filepath) //创建文件
		if err != nil {
			http.Error(w, "无法创建文件", http.StatusInternalServerError)
			errorLog.Println("无法创建文件", err)
			return
		}
		defer outfile.Close()

		_, err = io.Copy(outfile, file)
		if err != nil {
			http.Error(w, "无法保存文件", http.StatusInternalServerError)
			errorLog.Println("无法保存文件", err)
			return
		}
		fileChan <- err
	}()
	var oldpassword string
	err = tx.QueryRow("SELECT password FROM users WHERE username = ?", user).Scan(&oldpassword)
	if err != nil {
		errorLog.Println("数据库错误", err)
		return
	}
	if password == "" {
		_, err = tx.Exec("UPDATE users SET name = ? , age = ? , birthday = ? , photo = ? WHERE username = ?", name, age, birthday, filepath, user)
		if err != nil {
			errorLog.Println("数据库错误", err)
			return
		}
	} else if password != "" && password == oldpassword && password != newpassword {
		_, err = tx.Exec("UPDATE users SET name = ? , password = ? , age = ? , birthday = ? , photo = ? WHERE username = ?", name, newpassword, age, birthday, filepath, user)
		if err != nil {
			errorLog.Println("数据库错误", err)
			return
		}

	}
	if err = tx.Commit(); err != nil {
		http.Error(w, "服务器错误", http.StatusInternalServerError)
		return
	}
	mu.Unlock()
}

// 重置用户密码
func ResetpasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		_, err := db.Exec("UPDATE users SET password = 123456 WHERE username = ? ", username)
		if err != nil {
			errorLog.Println("数据库错误:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

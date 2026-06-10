package mode

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type Adjustbook struct {
	Adjust_id      int    `json:"adjust_id"`
	Adjust_date    string `json:"adjust_date"`
	Adjust_title   string `json:"adjust_title"`
	Adjust_isbn    string `json:"adjust_isbn"`
	Adjust_content string `json:"adjust_content"`
}

// LendBookHandler 处理借书请求（POST）和显示借书页面（GET）
func LendBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "html/lend-book-list.html", nil)
		return
	}

	// POST 请求需要身份验证（通过JWT）
	username := GetUsername(r)
	if username == "" {
		http.Error(w, "未登录", http.StatusUnauthorized)
		return
	}

	isbn := r.FormValue("isbn")
	expReturnDate := r.FormValue("exp_return_date")

	if isbn == "" || expReturnDate == "" {
		http.Error(w, "参数不完整", http.StatusBadRequest)
		return
	}

	// 验证日期格式
	_, err := time.Parse("2006-01-02", expReturnDate)
	if err != nil {
		http.Error(w, "日期格式无效", http.StatusBadRequest)
		return
	}

	// 检查是否已经借过同一本书
	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM cur_lend_records WHERE username = ? AND isbn = ?", username, isbn).Scan(&exists)
	if err != nil {
		errorLog.Println("查询借阅记录失败:", err)
		http.Error(w, "服务器错误", http.StatusInternalServerError)
		return
	}
	if exists > 0 {
		http.Error(w, "已借阅此书", http.StatusForbidden)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "服务器错误", http.StatusInternalServerError)
		return
	}

	// 获取图书信息和当前可借数量
	var curLendAmount, userCurLend int
	var title string
	err = tx.QueryRow("SELECT cur_lend_amount, title FROM all_books WHERE isbn = ?", isbn).Scan(&curLendAmount, &title)
	if err != nil {
		errorLog.Println("查询图书失败:", err)
		tx.Rollback()
		http.Error(w, "图书不存在", http.StatusNotFound)
		return
	}

	// 获取用户当前借书数量
	err = tx.QueryRow("SELECT user_cur_lend_amount FROM users WHERE username = ?", username).Scan(&userCurLend)
	if err != nil {
		errorLog.Println("查询用户失败:", err)
		tx.Rollback()
		http.Error(w, "用户不存在", http.StatusNotFound)
		return
	}

	// 检查是否达到借书上限（假设最多5本）
	if userCurLend >= 5 {
		tx.Rollback()
		http.Error(w, "已达借阅上限", http.StatusForbidden)
		return
	}

	// 检查是否有可借数量
	if curLendAmount <= 0 {
		tx.Rollback()
		http.Error(w, "此书已被借完", http.StatusNotFound)
		return
	}

	// 插入借阅记录
	lendDate := time.Now().Format("2006-01-02")
	_, err = tx.Exec("INSERT INTO lend_records (username, title, isbn, lend_date, exp_return_date) VALUES (?, ?, ?, ?, ?)",
		username, title, isbn, lendDate, expReturnDate)
	if err != nil {
		errorLog.Println("插入借阅记录失败:", err)
		tx.Rollback()
		http.Error(w, "借阅失败", http.StatusInternalServerError)
		return
	}

	// 插入当前借阅记录
	_, err = tx.Exec("INSERT INTO cur_lend_records (username, title, isbn, lend_date, exp_return_date) VALUES (?, ?, ?, ?, ?)",
		username, title, isbn, lendDate, expReturnDate)
	if err != nil {
		errorLog.Println("插入当前借阅记录失败:", err)
		tx.Rollback()
		http.Error(w, "借阅失败", http.StatusInternalServerError)
		return
	}

	// 更新图书可借数量
	_, err = tx.Exec("UPDATE all_books SET cur_lend_amount = cur_lend_amount - 1 WHERE isbn = ?", isbn)
	if err != nil {
		errorLog.Println("更新图书数量失败:", err)
		tx.Rollback()
		http.Error(w, "借阅失败", http.StatusInternalServerError)
		return
	}

	// 更新用户借阅数量
	_, err = tx.Exec("UPDATE users SET user_cur_lend_amount = user_cur_lend_amount + 1, user_his_lend_amount = user_his_lend_amount + 1 WHERE username = ?", username)
	if err != nil {
		errorLog.Println("更新用户数量失败:", err)
		tx.Rollback()
		http.Error(w, "借阅失败", http.StatusInternalServerError)
		return
	}

	// 更新汇总表
	_, err = tx.Exec("UPDATE library_summary SET total_lend_amount = total_lend_amount + 1")
	if err != nil {
		errorLog.Println("更新汇总表失败:", err)
		tx.Rollback()
		http.Error(w, "借阅失败", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		errorLog.Println("提交事务失败:", err)
		http.Error(w, "借阅失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "借阅成功"})
}

// ReturnBookHandler 处理还书请求（POST）
func ReturnBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := GetUsername(r)
		if username == "" {
			http.Error(w, "未登录", http.StatusUnauthorized)
			return
		}

		isbn := r.FormValue("isbn")
		if isbn == "" {
			http.Error(w, "ISBN不能为空", http.StatusBadRequest)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}

		// 检查用户是否有该书的借阅记录
		var title, lendDate, expReturnDate string
		err = tx.QueryRow("SELECT title, lend_date, exp_return_date FROM cur_lend_records WHERE username = ? AND isbn = ?", username, isbn).Scan(&title, &lendDate, &expReturnDate)
		if err != nil {
			errorLog.Println("查询借阅记录失败:", err)
			tx.Rollback()
			http.Error(w, "未找到借阅记录", http.StatusNotFound)
			return
		}

		returnDate := time.Now().Format("2006-01-02")

		// 计算逾期费用（每天0.1元）
		expDate, _ := time.Parse("2006-01-02", expReturnDate)
		retDate, _ := time.Parse("2006-01-02", returnDate)
		lateFee := 0.0
		if retDate.After(expDate) {
			days := int(retDate.Sub(expDate).Hours() / 24)
			lateFee = float64(days) * 0.1
		}

		// 插入还书记录
		_, err = tx.Exec("INSERT INTO return_records (username, title, isbn, lend_date, exp_return_date, return_date, late_fee) VALUES (?, ?, ?, ?, ?, ?, ?)",
			username, title, isbn, lendDate, expReturnDate, returnDate, lateFee)
		if err != nil {
			errorLog.Println("插入还书记录失败:", err)
			tx.Rollback()
			http.Error(w, "还书失败", http.StatusInternalServerError)
			return
		}

		// 删除当前借阅记录
		_, err = tx.Exec("DELETE FROM cur_lend_records WHERE username = ? AND isbn = ?", username, isbn)
		if err != nil {
			errorLog.Println("删除借阅记录失败:", err)
			tx.Rollback()
			http.Error(w, "还书失败", http.StatusInternalServerError)
			return
		}

		// 更新图书可借数量
		_, err = tx.Exec("UPDATE all_books SET cur_lend_amount = cur_lend_amount + 1 WHERE isbn = ?", isbn)
		if err != nil {
			errorLog.Println("更新图书数量失败:", err)
			tx.Rollback()
			http.Error(w, "还书失败", http.StatusInternalServerError)
			return
		}

		// 更新用户借阅数量
		_, err = tx.Exec("UPDATE users SET user_cur_lend_amount = user_cur_lend_amount - 1 WHERE username = ?", username)
		if err != nil {
			errorLog.Println("更新用户数量失败:", err)
			tx.Rollback()
			http.Error(w, "还书失败", http.StatusInternalServerError)
			return
		}

		// 更新汇总表
		_, err = tx.Exec("UPDATE library_summary SET total_return_amount = total_return_amount + 1, total_lend_amount = total_lend_amount - 1")
		if err != nil {
			errorLog.Println("更新汇总表失败:", err)
			tx.Rollback()
			http.Error(w, "还书失败", http.StatusInternalServerError)
			return
		}

		if err = tx.Commit(); err != nil {
			errorLog.Println("提交事务失败:", err)
			http.Error(w, "还书失败", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "还书成功",
			"late_fee": lateFee,
		})
	}
}

func ViewSearchBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// GET 请求也返回JSON数据（供前端初始化加载）
		handleJSONRequest(w, r)
		return
	}
	if r.Method == http.MethodPost {
		handleJSONRequest(w, r)
	}
}

func handleJSONRequest(w http.ResponseWriter, r *http.Request) {
	selsearch := r.URL.Query().Get("selsearch")
	inpsearch := r.URL.Query().Get("inpsearch")

	var rows *sql.Rows
	var err error

	if selsearch == "" && inpsearch == "" {
		rows, err = db.Query("SELECT title,author,isbn,press,press_date,price,cur_lend_amount,intro,cover FROM all_books")
	} else {
		search := "%" + inpsearch + "%"

		switch selsearch {
		case "1":
			rows, err = db.Query("SELECT title,author,isbn,press,press_date,price,cur_lend_amount,intro,cover FROM all_books WHERE author like ? or title like ? or isbn like ? or press like ? or press_date like ? or book_type like ?", search, search, search, search, search, search)
		case "2":
			rows, err = db.Query("SELECT title,author,isbn,press,press_date,price,cur_lend_amount,intro,cover FROM all_books WHERE title like ?", search)
		case "3":
			rows, err = db.Query("SELECT title,author,isbn,press,press_date,price,cur_lend_amount,intro,cover FROM all_books WHERE author like ?", search)
		case "4":
			rows, err = db.Query("SELECT title,author,isbn,press,press_date,price,cur_lend_amount,intro,cover FROM all_books WHERE isbn like ?", search)
		default:
			http.Error(w, "无效的搜索类型", http.StatusBadRequest)
			return
		}
	}

	if err != nil {
		errorLog.Println("数据库错误：", err)
		http.Error(w, "服务器错误", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Books
	for rows.Next() {
		var Book Books
		err = rows.Scan(&Book.Title, &Book.Author, &Book.ISBN, &Book.Press, &Book.Press_date, &Book.Price, &Book.Cur_Lend_amount, &Book.Intro, &Book.Cover)
		if err != nil {
			errorLog.Println("数据库错误:", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}
		err = db.QueryRow("SELECT rec_type FROM recommend_books WHERE isbn = ?", Book.ISBN).Scan(&Book.Rec_type)
		if err != nil && err != sql.ErrNoRows {
			errorLog.Println("获取推荐类型错误:", err)
		}
		books = append(books, Book)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		errorLog.Println("编码错误：", err)
		http.Error(w, "服务器错误", http.StatusInternalServerError)
		return
	}
}

func ViewAdjustBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "html/view-adjustbook.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		rows, err := db.Query("SELECT * FROM adjust_books")
		if err != nil {
			errorLog.Println("数据库错误", err)
			return
		}
		defer rows.Close()

		var adjusts []Adjustbook
		for rows.Next() {
			var adjust Adjustbook
			err = rows.Scan(&adjust.Adjust_id, &adjust.Adjust_date, &adjust.Adjust_title, &adjust.Adjust_isbn, &adjust.Adjust_content)
			if err != nil {
				errorLog.Println("数据库错误：", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			adjusts = append(adjusts, adjust)
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(adjusts)
		if err != nil {
			errorLog.Println("编码错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func ClassifySearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var rows *sql.Rows
		var err error

		value := r.FormValue("value")
		value1 := r.FormValue("value1")
		value2 := r.FormValue("value2")
		value4 := r.FormValue("value4")
		date := r.FormValue("value3")

		value3 := "%" + date + "%"

		query := "SELECT title, author, isbn, press, press_date, price, cur_lend_amount, intro, cover FROM all_books WHERE 1=1 "

		var params []interface{}

		if value != "" {
			query += " AND author = ?"
			params = append(params, value)
		}
		if value1 != "" {
			query += " AND press = ?"
			params = append(params, value1)
		}
		if value2 != "" {
			query += " AND book_type = ?"
			params = append(params, value2)
		}
		if value3 != "" && len(date) > 0 {
			query += " AND press_date LIKE ?"
			params = append(params, value3)
		}
		if value4 != "" {
			a := 1
			query += " AND rec_state = ?"
			params = append(params, a)
		}

		// 执行查询
		rows, err = db.Query(query, params...)
		if err != nil {
			errorLog.Println("数据库错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer rows.Close()
		var books []Books
		for rows.Next() {
			var Book Books
			err = rows.Scan(&Book.Title, &Book.Author, &Book.ISBN, &Book.Press, &Book.Press_date, &Book.Price, &Book.Cur_Lend_amount, &Book.Intro, &Book.Cover)
			if err != nil {
				errorLog.Println("数据库错误:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			err = db.QueryRow("SELECT rec_type FROM recommend_books WHERE isbn = ?", Book.ISBN).Scan(&Book.Rec_type)
			if err != nil && err == sql.ErrNoRows {
				// rec_type 为空是正常情况
			}
			books = append(books, Book)
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(books)
		if err != nil {
			errorLog.Println("编码错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
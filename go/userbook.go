package mode

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
	// "golang.org/x/text/date"
	// "golang.org/x/tools/go/analysis/passes/defers"
)

func LendBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("html/lend-book-list.html")
		if err != nil {
			fmt.Printf("解析模板失败: %v\n", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
		}
		err = tmpl.ExecuteTemplate(w, "lend-book-list.html", nil)
		if err != nil {
			fmt.Printf("执行模板失败: %v\n", err)
			errorLog.Println("服务器错误：", err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
		}
	}
	if r.Method == http.MethodPost {
		var rows *sql.Rows
		var err error

		category := r.FormValue("category")
		value := r.FormValue("value")

		switch category {
		case "作者":
			rows, err = db.Query("SELECT title,author,isbn,press,press_date,price,cur_lend_amount,intro,cover FROM all_books WHERE author = ?", value)
		case "出版社":
			rows, err = db.Query("SELECT title,author,isbn,press,press_date,price,cur_lend_amount,intro,cover FROM all_books WHERE press = ?", value)
		case "出版日期":
			values := value + "%"
			rows, err = db.Query("SELECT title,author,isbn,press,press_date,price,cur_lend_amount,intro,cover FROM all_books WHERE press_date LIKE ?", values)
		case "类型":
			rows, err = db.Query("SELECT title,author,isbn,press,press_date,price,cur_lend_amount,intro,cover FROM all_books WHERE book_type = ?", value)
		default:
			rows, err = db.Query("SELECT title,author,isbn,press,press_date,price,cur_lend_amount,intro,cover FROM all_books")
		}
		if err != nil {
			errorLog.Println("数据库错误:", err)
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

		session, _ := store.Get(r, "user-session")
		user, ok := session.Values["username"].(string)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

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

		isbn := r.FormValue("isbn")

		var book Books
		err = tx.QueryRow("SELECT title,isbn,cur_lend_amount,rec_state FROM all_books WHERE isbn = ?", isbn).Scan(&book.Title, &book.ISBN, &book.Cur_Lend_amount, &book.Rec_state)

		if err != nil {
			errorLog.Println("数据库错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var cur int
		err = tx.QueryRow("SELECT user_cur_lend_amount from users WHERE username = ?", user).Scan(&cur)
		if err != nil {
			errorLog.Println("数据库错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if cur > 5 {
			errorLog.Println("当前借书已达上限")
			w.WriteHeader(http.StatusForbidden) //403
			return
		}
		var id int
		err = tx.QueryRow("SELECT lend_id FROM cur_lend_records WHERE username = ? AND isbn = ?", user, isbn).Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				id = 1
			} else {
				errorLog.Println("数据库错误：", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if id == 1 {
			if book.Cur_Lend_amount > 0 {
				_, err := tx.Exec("UPDATE all_books set cur_lend_amount = cur_lend_amount -1 WHERE isbn = ?", isbn)
				if err != nil {
					errorLog.Println("数据库错误：", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				_, err = tx.Exec("UPDATE library_summary set total_lend_amount = total_lend_amount -1 , total_return_amount = total_return_amount +1")
				if err != nil {
					errorLog.Println("数据库错误：", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				_, err = tx.Exec("UPDATE users SET user_cur_lend_amount = user_cur_lend_amount + 1 , user_his_lend_amount = user_his_lend_amount +1 WHERE username = ?", user)
				if err != nil {
					errorLog.Println("数据库错误：", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				currentTime := time.Now()
				now := currentTime.Format("2006-01-02")
				exp_return_date := r.FormValue("exp_return_date")

				expTime := time.Now().Truncate(24 * time.Hour)
				expreturn, err := time.Parse("2006-01-02", exp_return_date)
				if err != nil {
					errorLog.Println("日期转换错误")
					return
				}
				// 计算两个日期之间的差值
				day := expreturn.Sub(expTime.Truncate(24*time.Hour)).Hours() / 24

				if day >= 0 && day <= 14 {
					_, err = tx.Exec("INSERT INTO lend_records (username,title,isbn,lend_date,exp_return_date) values(?,?,?,?,?)", user, book.Title, book.ISBN, now, exp_return_date)
					if err != nil {
						errorLog.Println("数据库错误：", err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				} else {
					errorLog.Println("日期设置错误")
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				_, err = tx.Exec("INSERT INTO cur_lend_records (username,title,isbn,lend_date,exp_return_date) values(?,?,?,?,?)", user, book.Title, book.ISBN, now, exp_return_date)
				if err != nil {
					errorLog.Println("数据库错误：", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				if book.Rec_state == 1 {
					_, err = tx.Exec("UPDATE recommend_books SET cur_lend_amount = cur_lend_amount - 1 WHERE isbn = ?", isbn)
					if err != nil {
						errorLog.Println("数据库错误：", err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
				if err = tx.Commit(); err != nil {
					http.Error(w, "服务器错误", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			errorLog.Println("此书已被借完")
			return
		}
		errorLog.Println("已借出这本书")
		w.WriteHeader(http.StatusForbidden) //403
		return

	}
}

func ReturnBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		session, _ := store.Get(r, "user-session")
		user, ok := session.Values["username"].(string)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

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

		isbn := r.FormValue("isbn")

		var expreturndate string
		var title string
		var lenddate string

		err = tx.QueryRow("SELECT title,lend_date,exp_return_date FROM cur_lend_records WHERE username = ? AND isbn = ? ", user, isbn).Scan(&title, &lenddate, &expreturndate)
		if err != nil {
			errorLog.Println("数据库错误:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		expReturnDate, err := time.Parse("2006-01-02", expreturndate) // 确保格式与数据库中的一致
		if err != nil {
			errorLog.Println("解析日期错误:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		currentTime := time.Now().Truncate(24 * time.Hour)

		// 计算两个日期之间的差值
		day := currentTime.Sub(expReturnDate.Truncate(24*time.Hour)).Hours() / 24

		var money float64

		var price float64

		if day > 0 {
			err = tx.QueryRow("SELECT price FROM all_books WHERE isbn = ?", isbn).Scan(&price)
			if err != nil {
				errorLog.Println("数据库错误:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			money = price * day * 0.02
		} else {
			money = 0
		}

		moneykey := r.FormValue("money-key")
		if money > 0 {
			if moneykey != "123456" {
				w.WriteHeader(http.StatusForbidden) //403
				return
			}
			w.WriteHeader(http.StatusOK)
		}

		_, err = tx.Exec("UPDATE all_books  SET cur_lend_amount = cur_lend_amount +1 WHERE isbn = ?", isbn)
		if err != nil {
			errorLog.Println("数据库错误:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("UPDATE users SET user_cur_lend_amount = user_cur_lend_amount -1 WHERE username = ?", user)
		if err != nil {
			errorLog.Println("数据库错误:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var rec_state int
		err = tx.QueryRow("SELECT rec_state FROM all_books WHERE isbn = ?", isbn).Scan(&rec_state)
		if err != nil {
			errorLog.Println("数据库错误:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if rec_state == 1 {
			_, err = tx.Exec("UPDATE recommend_books SET cur_lend_amount = cur_lend_amount + 1 WHERE isbn = ?", isbn)
			if err != nil {
				errorLog.Println("数据库错误:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		_, err = tx.Exec("UPDATE library_summary set total_lend_amount = total_lend_amount +1 , total_return_amount = total_return_amount -1")
		if err != nil {
			errorLog.Println("数据库错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		Now := time.Now()
		now := Now.Format("2006-01-02")

		_, err = tx.Exec("INSERT INTO return_records (username,title,isbn,lend_date,exp_return_date,return_date,late_fee) values(?,?,?,?,?,?,?)", user, title, isbn, lenddate, expreturndate, now, money)
		if err != nil {
			errorLog.Println("数据库错误:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("DELETE FROM cur_lend_records WHERE username = ? AND isbn = ?", user, isbn)
		if err != nil {
			errorLog.Println("数据库错误:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = tx.Commit(); err != nil {
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func ViewSearchBookHandler(w http.ResponseWriter, r *http.Request) {
	var rows *sql.Rows
	var err error

	selsearch := r.FormValue("selsearch")
	inpsearch := r.FormValue("inpsearch")

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
		w.WriteHeader(http.StatusBadRequest) //有错误返回错误状态码400
		return
	}
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

func ViewOnlyBookHandler(w http.ResponseWriter, r *http.Request) {

}

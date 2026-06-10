package mode

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

)

// 所有图书结构体
type Books struct {
	Title           string  `json:"title"`
	Author          string  `json:"author"`
	Price           float64 `json:"price"`
	Rec_state       int     `json:"rec_state"`
	Press           string  `json:"press"`
	ISBN            string  `json:"isbn"`
	Intro           string  `json:"intro"`
	Amount          int     `json:"amount"`
	Lend_amount     int     `json:"lend_amount"`
	Cur_Lend_amount int     `json:"cur_lend_amount"`
	Book_type       string  `json:"book_type"`
	Press_date      string  `json:"press_date"`
	Cover           string  `json:"cover"`
	Rec_type        string  `json:"rec_type"`
}

type LendRecords struct {
	Lend_id         int    `json:"lend_id"`
	Username        string `json:"username"`
	Title           string `json:"title"`
	ISBN            string `json:"isbn"`
	Lend_date       string `json:"lend_date"`
	Exp_return_date string `json:"exp_return_date"`
}

type ReturnRecords struct {
	Return_id       int     `json:"return_id"`
	Username        string  `json:"username"`
	Title           string  `json:"title"`
	ISBN            string  `json:"isbn"`
	Lend_date       string  `json:"lend_date"`
	Exp_return_date string  `json:"exp_return_date"`
	Return_date     string  `json:"return_date"`
	Late_fee        float64 `json:"late_fee"`
}

// 添加图书功能
func AddBooksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "html/addbook.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseMultipartForm(10 << 20)
		ctx := r.Context()
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			errorLog.Println("无法处理事务", err)
			return
		}
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
		title := r.FormValue("title")
		author := r.FormValue("author")
		book_type := r.FormValue("book_type")
		press := r.FormValue("press")
		press_date := r.FormValue("press_date")
		isbn := r.FormValue("isbn")
		intro := r.FormValue("intro")
		price := r.FormValue("price")
		amount := r.FormValue("amount")
		lend_amount := r.FormValue("lend_amount")
		rec_state := r.FormValue("rec_state")

		file, _, err := r.FormFile("cover") //获取文件
		if err != nil {
			http.Error(w, "无法获取文件", http.StatusBadRequest)
			errorLog.Println("无法获取文件", err)
			return
		}
		defer file.Close()
		filename := "image" + isbn + ".jpg"           //文件命名
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

		var bookISBN string
		err = db.QueryRow("SELECT ISBN FROM all_books WHERE ISBN = ?", isbn).Scan(&bookISBN) //查询输入的isbn是否存在，查询结果赋给变量
		if err != sql.ErrNoRows {                                                            //如果不为空
			if err != nil { //其他错误
				w.WriteHeader(http.StatusInternalServerError)
				errorLog.Println("查询失败：", err)
				return
			}
			errorLog.Println("已有这个isbn")
			w.WriteHeader(http.StatusConflict) //已存在返回409
			return
		}

		if rec_state == "1" { //如果这本书是推荐书籍
			rec_type := r.FormValue("rec_type") //获取一个推荐类型文本
			if rec_type == "" {                 //如果没有写
				rec_type = book_type //推荐类型就=书本类型
			}
			//添加到推荐图书表
			_, err = tx.Exec("INSERT INTO recommend_books (isbn,title,author,rec_type,cover,cur_lend_amount) values(?,?,?,?,?,?)", isbn, title, author, rec_type, filepath, lend_amount)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				errorLog.Println("数据库错误：", err)
				return
			}
		} else {
			rec_state = "0" //不推荐传入0
		}

		if amount >= lend_amount { //判断总数量是否大于等于可借数量
			//添加到所有图书表
			_, err = tx.Exec("INSERT INTO all_books (title,author,book_type,press,press_date,isbn,intro,price,amount,lend_amount,cur_lend_amount,rec_state,cover) values(?,?,?,?,?,?,?,?,?,?,?,?,?)", title, author, book_type, press, press_date, isbn, intro, price, amount, lend_amount, lend_amount, rec_state, filepath)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				errorLog.Println("数据库错误：", err)
				return
			}
			//更新汇总表信息
			_, err = tx.Exec("UPDATE library_summary SET total_books_amount = total_books_amount + ? , total_lend_amount = total_lend_amount + ?", amount, lend_amount)
			if err != nil {
				errorLog.Println("数据库错误：", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity) //数据错误返回状态码422
			return
		}
		if err := <-fileChan; err != nil {
			errorLog.Println("文件上传失败", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			errorLog.Println("无法提交事务", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// 查询书籍功能
func ViewBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "html/view-book.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		var row *sql.Rows
		// var rows *sql.Rows
		var err error
		class := r.FormValue("category") //获取select标签的value

		stitle := "%" + title + "%" //用作模糊查询

		switch class {
		case "1": //value = 1
			row, err = db.Query("SELECT * FROM all_books WHERE title LIKE ? OR author LIKE ? OR book_type LIKE ? OR press LIKE ? OR press_date LIKE ? OR isbn LIKE ?", stitle, stitle, stitle, stitle, stitle, stitle) //任意词查询

		case "2":
			row, err = db.Query("SELECT * FROM all_books WHERE title = ?", title) //标题查询

		case "3":
			row, err = db.Query("SELECT * FROM all_books WHERE author = ?", title) //作者查询

		case "4":
			row, err = db.Query("SELECT * FROM all_books WHERE isbn = ?", title) //isbn查询

		default:
			w.WriteHeader(http.StatusBadRequest) //有错误返回错误状态码400
			return
		}
		if err != nil {
			if err == sql.ErrNoRows {
				errorLog.Println("没有这本书：", err)
				w.WriteHeader(http.StatusNotFound)
			} else {
				errorLog.Println("数据库错误：", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		defer row.Close()

		var books []Books
		for row.Next() {
			var book Books
			err := row.Scan(&book.Title, &book.Author, &book.Book_type, &book.Press, &book.Press_date, &book.ISBN, &book.Cover, &book.Intro, &book.Price, &book.Amount, &book.Lend_amount, &book.Cur_Lend_amount, &book.Rec_state)
			if err != nil {
				errorLog.Println("服务器错误：", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			db.QueryRow("SELECT rec_type FROM recommend_books WHERE isbn = ?", book.ISBN).Scan(&book.Rec_type)

			books = append(books, book)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(books)
		if err != nil {
			errorLog.Println("编码错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
}

// 更新图书信息
func UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// yisbn := "1"
		yisbn := r.FormValue("yisbn") //获取原图书的isbn
		title := r.FormValue("title")
		author := r.FormValue("author")
		book_type := r.FormValue("book_type")
		press := r.FormValue("press")
		press_date := r.FormValue("press_date")
		isbn := r.FormValue("isbn")
		intro := r.FormValue("intro")
		price := r.FormValue("price")
		amount := r.FormValue("amount")
		lend_amount := r.FormValue("lend_amount")
		cur_lend_amount := r.FormValue("cur_lend_amount")
		rec_state := r.FormValue("rec_state")

		num, err := strconv.Atoi(amount) //将接收到的数据转换成数值型
		if err != nil {
			errorLog.Println("数量转换错误:", err)
			return
		}

		num1, err := strconv.Atoi(lend_amount)
		if err != nil {
			errorLog.Println("数量转换错误:", err)
			return
		}

		num2, err := strconv.Atoi(cur_lend_amount)
		if err != nil {
			errorLog.Println("数量转换错误:", err)
			return
		}

		// 封面可选：如果没有上传新封面，保留原有的封面路径
		var coverPath string
		file, _, err := r.FormFile("cover")
		if err == nil {
			defer file.Close()
			filename := "image" + isbn + ".jpg"
			coverPath = filepath.Join("images", filename)

			outfile, err := os.Create(coverPath)
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
		} else {
			// 没有上传封面，获取原有的封面路径
			err = db.QueryRow("SELECT cover FROM all_books WHERE isbn = ?", yisbn).Scan(&coverPath)
			if err != nil {
				coverPath = ""
			}
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

		var bookisbn string
		if isbn != yisbn { //如果更新图书的isbn不等于原isbn
			err = tx.QueryRow("SELECT ISBN FROM all_books WHERE isbn = ?", isbn).Scan(&bookisbn) //在总的图书表中查询输入的isbn，是否存在
			if err != sql.ErrNoRows {                                                            //存在
				if err != nil { //其他错误
					http.Error(w, "查询失败", http.StatusInternalServerError)
					errorLog.Println("查询失败：", err)
					return
				}
				errorLog.Println("已有这个ISBN")       //返回错误信息
				w.WriteHeader(http.StatusConflict) //409
				return                             // 退出程序
			}
		}
		// isbn等于原isbn或不存在
		var book_amount, book_lend_amount, book_cur_lend_amount int //定义原数数量，可借数量，当前可借数量
		row := tx.QueryRow("SELECT amount, lend_amount, cur_lend_amount FROM all_books WHERE isbn = ?", yisbn)
		err = row.Scan(&book_amount, &book_lend_amount, &book_cur_lend_amount) //获取原书的值赋给变量
		//判断可借数量和当前可借数量的差是否一致，书本数量是否大于等于可借数量，可借数量是否大于等于当前可借数量
		if (book_lend_amount-book_cur_lend_amount) == (num1-num2) && amount >= lend_amount && lend_amount >= cur_lend_amount {
			rec_type := r.FormValue("rec_type") //获取推荐类型

			if rec_state == "1" { //如果推荐
				row := tx.QueryRow("SELECT isbn FROM recommend_books WHERE isbn = ?", yisbn) //查询原isbn是否在推荐图书表中
				err := row.Scan(&bookisbn)
				if rec_type == "" {
					rec_type = book_type
				}
				if err == sql.ErrNoRows { //不存在
					//新插入一条记录到推荐图书表
					_, err = tx.Exec("INSERT INTO recommend_books (isbn,title,author,rec_type,cover,cur_lend_amount) values(?,?,?,?,?,?)", isbn, title, author, rec_type, coverPath, cur_lend_amount)
					if err != nil {
						errorLog.Println("数据库错误：", err)
						return
					}
				} else if err != nil {
					errorLog.Println("数据库错误：", err)
					return
				} else { //存在
					//更新这条记录
					_, err = tx.Exec("UPDATE recommend_books SET isbn = ? , title = ? , author = ? , rec_type = ? , cover = ? , cur_lend_amount = ? WHERE isbn = ?", isbn, title, author, rec_type, coverPath, cur_lend_amount, yisbn)
					if err != nil {
						errorLog.Println("数据库错误：", err)
						return
					}
				}
			} else {
				rec_state = "0"                                                       //不推荐
				_, err = tx.Exec("DELETE FROM recommend_books WHERE isbn = ?", yisbn) //删除推荐图书表的这条记录
				if err != nil {
					errorLog.Println("数据库错误：", err)
					return
				}
			}
			//根据原isbn更新所有图书表的信息
			_, err = tx.Exec("UPDATE all_books SET title = ? , author = ? ,book_type = ? , press = ? , press_date = ? , isbn = ? , intro = ? , price = ? , amount = ? , lend_amount = ?, cur_lend_amount = ? , rec_state = ? , cover = ? WHERE isbn = ?", title, author, book_type, press, press_date, isbn, intro, price, amount, lend_amount, cur_lend_amount, rec_state, coverPath, yisbn)
			if err != nil {
				errorLog.Println("数据库错误：", err)
				return
			}

			_, err = tx.Exec("UPDATE cur_lend_records SET title = ? ,isbn = ? WHERE isbn = ?", title, isbn, yisbn)
			if err != nil {
				errorLog.Println("数据库错误：", err)
				return
			}

			_, err = tx.Exec("UPDATE lend_records SET title = ? ,isbn = ? WHERE isbn = ?", title, isbn, yisbn)
			if err != nil {
				errorLog.Println("数据库错误：", err)
				return
			}

			_, err = tx.Exec("UPDATE return_records SET title = ? , isbn = ? WHERE isbn = ?", title, isbn, yisbn)
			if err != nil {
				errorLog.Println("数据库错误：", err)
				return
			}

			//同步更新汇总表的内容
			_, err = tx.Exec("UPDATE library_summary SET total_books_amount = total_books_amount + ? , total_lend_amount = total_lend_amount + ?", num-book_amount, num2-book_cur_lend_amount)
			if err != nil {
				errorLog.Println("数据库错误：", err)
				return
			}
			if err = tx.Commit(); err != nil {
				http.Error(w, "服务器错误", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		errorLog.Println("数量设置错误")
		w.WriteHeader(http.StatusUnprocessableEntity) //422
	}
}

// 删除图书
func DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// yisbn := "5"
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
		yisbn := r.FormValue("yisbn") //获取图书isbn
		var num1, num2 int
		var rec_state string
		row := tx.QueryRow("SELECT lend_amount,cur_lend_amount,rec_state FROM all_books WHERE isbn = ?", yisbn)
		err = row.Scan(&num1, &num2, &rec_state) //获取图书可借数量，当前可借数量，推荐状态
		if err != nil {
			errorLog.Println("查询错误", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if num1 == num2 { //可借数量和当前可借数量一样（没有书借出）
			//执行删除操作
			_, err = tx.Exec("DELETE FROM all_books WHERE isbn = ?", yisbn)
			if err != nil {
				errorLog.Println("服务器错误", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			//如果是推荐图书
			if rec_state == "1" {
				//同步删除推荐图书表中的记录
				_, err = tx.Exec("DELETE FROM recommend_books WHERE isbn = ?", yisbn)
				if err != nil {
					errorLog.Println("服务器错误", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			//删除完同步更新汇总表的内容
			_, err = tx.Exec("UPDATE library_summary SET total_books_amount = total_books_amount - ? , total_lend_amount = total_lend_amount - ?", num1, num2)
			if err != nil {
				errorLog.Println("服务器错误", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if err = tx.Commit(); err != nil {
				http.Error(w, "服务器错误", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else {
			//还有未归还的书
			errorLog.Println("此书未全部归还，无法删除")        //返回错误信息
			w.WriteHeader(http.StatusNotAcceptable) //返回错误状态吗406
		}
	}
}

func AdjustBookHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "html/adjust-book.html", nil)


	if r.Method == http.MethodPost {
		adjust_date := r.FormValue("adjust_date")
		adjust_title := r.FormValue("adjust_title")
		adjust_isbn := r.FormValue("adjust_isbn")
		adjust_content := r.FormValue("adjust_content")

		adjustTime := time.Now().Truncate(24 * time.Hour)
		adjustdate, _ := time.Parse("2006-01-02", adjust_date)
		// 计算两个日期之间的差值
		day := adjustdate.Sub(adjustTime.Truncate(24*time.Hour)).Hours() / 24
		if day == 0 {
			_, err := db.Exec("INSERT INTO adjust_books (adjust_date,adjust_title,adjust_isbn,adjust_content) values(?,?,?,?)", adjust_date, adjust_title, adjust_isbn, adjust_content)
			if err != nil {
				errorLog.Println("数据库错误：", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}

}

func ViewLendRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "html/view-lend-records.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		var rows *sql.Rows
		var err error
		class := r.FormValue("class")
		lend_date := r.FormValue("lend_date")

		switch class {
		case "1":
			rows, err = db.Query("SELECT * FROM lend_records")
		case "2":
			rows, err = db.Query("SELECT * FROM lend_records WHERE lend_date = ?", lend_date)
		case "3":
			rows, err = db.Query("SELECT * FROM lend_records WHERE username = ?", lend_date)
		default:
			w.WriteHeader(http.StatusBadRequest) //有错误返回错误状态码400
			return
		}
		if err != nil {
			errorLog.Println("数据库错误", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var lendrecords []LendRecords
		for rows.Next() {
			var lendrecord LendRecords
			err = rows.Scan(&lendrecord.Lend_id, &lendrecord.Username, &lendrecord.Title, &lendrecord.ISBN, &lendrecord.Lend_date, &lendrecord.Exp_return_date)
			if err != nil {
				errorLog.Println("数据库错误", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			lendrecords = append(lendrecords, lendrecord)
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(lendrecords)
		if err != nil {
			errorLog.Println("编码错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return

	}
}

func ViewReturnRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "html/view-return-records.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		var rows *sql.Rows
		var err error
		class := r.FormValue("class")
		return_date := r.FormValue("return_date")

		switch class {
		case "1":
			rows, err = db.Query("SELECT * FROM return_records")
		case "2":
			rows, err = db.Query("SELECT * FROM return_records WHERE return_date = ?", return_date)
		case "3":
			rows, err = db.Query("SELECT * FROM lend_records WHERE username = ?", return_date)
		default:
			w.WriteHeader(http.StatusBadRequest) //有错误返回错误状态码400
			return
		}

		if err != nil {
			errorLog.Println("数据库错误", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var returnrecords []ReturnRecords
		for rows.Next() {
			var returnrecord ReturnRecords
			err = rows.Scan(&returnrecord.Return_id, &returnrecord.Username, &returnrecord.Title, &returnrecord.ISBN, &returnrecord.Lend_date, &returnrecord.Exp_return_date, &returnrecord.Return_date, &returnrecord.Late_fee)
			if err != nil {
				errorLog.Println("数据库错误", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			returnrecords = append(returnrecords, returnrecord)
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(returnrecords)
		if err != nil {
			errorLog.Println("编码错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return

	}
}

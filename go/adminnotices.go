package mode

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"text/template"
	"time"
)

type Notices struct {
	Notice_id    int    `json:"notice_id"`
	Notice_date  string `json:"notice_date"`
	Notice_title string `json:"notice_title"`
	Notice       string `json:"notice"`
}

// 新增公告
func AddNoticeHandler(w http.ResponseWriter, r *http.Request) {
	//打开网页
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
			tmpl, err := template.ParseFiles("html/addnotice.html")
			if err != nil {
				fmt.Printf("解析模板失败: %v\n", err)
				http.Error(w, "服务器错误", http.StatusInternalServerError)
			}
			err = tmpl.ExecuteTemplate(w, "addnotice.html", nil)
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
		notice_date := r.FormValue("notice_date")
		notice_title := r.FormValue("notice_title")
		notice := r.FormValue("notice")
		//插入公告记录
		noticeTime := time.Now().Truncate(24 * time.Hour)
		noticedate, _ := time.Parse("2006-01-02", notice_date)
		// 计算两个日期之间的差值
		day := noticedate.Sub(noticeTime.Truncate(24*time.Hour)).Hours() / 24
		if day >= 0 {
			_, err := db.Exec("INSERT INTO notices (notice_date,notice_title,notice) values(?,?,?)", notice_date, notice_title, notice)
			if err != nil {
				errorLog.Println("数据库错误", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest) //400
		}
	}
}

// 查询公告内容
func ViewNoticeHandler(w http.ResponseWriter, r *http.Request) {
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
			tmpl, err := template.ParseFiles("html/view-notice.html")
			if err != nil {
				fmt.Printf("解析模板失败: %v\n", err)
				http.Error(w, "服务器错误", http.StatusInternalServerError)
			}
			err = tmpl.ExecuteTemplate(w, "view-notice.html", nil)
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
		notice_id := r.FormValue("notice_id") //获取查询内容
		chaxtype := r.FormValue("chaxtype")   //获取查询类型

		var row *sql.Rows
		var err error

		allnotice := "%" + notice_id + "%" //模糊查询

		switch chaxtype {
		case "0":
			//任意词匹配
			row, err = db.Query("SELECT * FROM notices WHERE notice_id LIKE ? or notice_date LIKE ? or notice_title LIKE ? or notice LIKE ?", allnotice, allnotice, allnotice, allnotice)
		case "1":
			//公告id匹配
			row, err = db.Query("SELECT * FROM notices WHERE notice_id = ?", notice_id)
		case "2":
			//公告日期匹配
			row, err = db.Query("SELECT * FROM notices WHERE notice_date = ?", notice_id)
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err != nil {
			if err == sql.ErrNoRows {
				errorLog.Println("没有这条公告：", err)
				w.WriteHeader(http.StatusNotFound)
			} else {
				errorLog.Println("数据库错误：", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		defer row.Close()

		var notices []Notices //定义切片
		for row.Next() {
			var notice Notices                                                                            //定义结构体
			err := row.Scan(&notice.Notice_id, &notice.Notice_date, &notice.Notice_title, &notice.Notice) //获取数据
			if err != nil {
				errorLog.Println("服务器错误：", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			notices = append(notices, notice) //循环累加到切片
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(notices)
		if err != nil {
			errorLog.Println("编码错误：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// 更新公告内容
func UpdateNoticeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		notice_id := r.FormValue("notice_id")       //公告id
		notice_date := r.FormValue("notice_date")   //公告日期
		notice_title := r.FormValue("notice_title") //公告标题
		notice := r.FormValue("notice")             //公告内容

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}

		var noticedate string
		err = db.QueryRow("SELECT notice_date FROM notices WHERE notice_id = ? ", notice_id).Scan(&noticedate)

		if err != nil {
			errorLog.Println("数据库错误")
			return
		}

		hisnoticedate, err := time.Parse("2006-01-02", noticedate) // 确保格式与数据库中的一致
		if err != nil {
			errorLog.Println("解析日期错误:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		nownoticedate, err := time.Parse("2006-01-02", notice_date) // 确保格式与数据库中的一致
		if err != nil {
			errorLog.Println("解析日期错误:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		now := time.Now().Truncate(24 * time.Hour)
		// 计算两个日期之间的差值
		// day := nownoticedate.Sub(hisnoticedate.Truncate(24*time.Hour)).Hours() / 24

		day := nownoticedate.Sub(now.Truncate(24*time.Hour)).Hours() / 24

		day2 := hisnoticedate.Sub(now.Truncate(24*time.Hour)).Hours() / 24

		if day >= 0 && day2 >= 0 {
			_, err = tx.Exec("UPDATE notices SET notice_date = ? , notice_title = ? , notice = ? WHERE notice_id = ?", notice_date, notice_title, notice, notice_id) //更新
			if err != nil {
				tx.Rollback()
				errorLog.Println("数据库错误", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if err = tx.Commit(); err != nil {
				http.Error(w, "服务器错误", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

// 删除公告
func DeleteNoticeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		notice_id := r.FormValue("notice_id")
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("DELETE FROM notices WHERE notice_id = ?", notice_id)
		if err != nil {
			tx.Rollback()
			errorLog.Println("数据库错误", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = tx.Commit(); err != nil {
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}
	}
}

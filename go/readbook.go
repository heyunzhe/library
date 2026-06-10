package mode

import (
	"encoding/json"
	"html/template"
	"net/http"
)

// UploadContentHandler 管理员上传图书内容
func UploadContentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		isbn := r.FormValue("isbn")
		content := r.FormValue("content")
		if isbn == "" || content == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "参数不完整"})
			return
		}

		_, err := db.Exec("REPLACE INTO book_contents (isbn, content) VALUES (?, ?)", isbn, content)
		if err != nil {
			errorLog.Println("保存图书内容失败:", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "保存失败"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "内容上传成功"})
		return
	}
	http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
}

// ReadBookHandler 在线阅读（仅已借阅用户可以看）
func ReadBookHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUsername(r)
	if user == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		isbn := r.URL.Query().Get("isbn")
		if isbn == "" {
			http.Error(w, "缺少参数", http.StatusBadRequest)
			return
		}

		// 检查是否借阅了这本书
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM cur_lend_records WHERE username = ? AND isbn = ?", user, isbn).Scan(&count)
		if err != nil || count == 0 {
			// 再查历史借阅
			err = db.QueryRow("SELECT COUNT(*) FROM lend_records WHERE username = ? AND isbn = ?", user, isbn).Scan(&count)
			if err != nil || count == 0 {
				http.Error(w, "您没有借阅这本书", http.StatusForbidden)
				return
			}
		}

		// 获取图书信息和内容
		var title, author, intro, content string
		err = db.QueryRow("SELECT title, author, COALESCE(intro,''), COALESCE((SELECT content FROM book_contents WHERE isbn=?),'') FROM all_books WHERE isbn=?", isbn, isbn).
			Scan(&title, &author, &intro, &content)
		if err != nil {
			http.Error(w, "图书不存在", http.StatusNotFound)
			return
		}

		if content == "" {
			content = "暂无内容，请等待管理员上传。"
		}

		tmpl, err := template.ParseFiles("html/read-book.html")
		if err != nil {
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, map[string]interface{}{
			"Title":   title,
			"Author":  author,
			"Isbn":    isbn,
			"Intro":   intro,
			"Content": content,
		})
		return
	}
	http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
}

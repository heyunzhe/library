package mode

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"sync"
	"text/template"
)

var errorLog *log.Logger
var db *sql.DB

// templateCache 模板缓存，避免每次请求重复解析
var templateCache sync.Map

// renderTemplate 从缓存加载并渲染模板
func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := templateCache.Load(name)
	if !ok {
		parsed, err := template.ParseFiles(name)
		if err != nil {
			errorLog.Printf("解析模板失败 %s: %v", name, err)
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}
		templateCache.Store(name, parsed)
		tmpl = parsed
	}
	err := tmpl.(*template.Template).Execute(w, data)
	if err != nil {
		errorLog.Printf("执行模板失败 %s: %v", name, err)
		http.Error(w, "服务器错误", http.StatusInternalServerError)
	}
}

func Init() {
	var err error
	//创建日志文件
	logfile, err := os.OpenFile("library.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("无法打开日志文件:", err)
	}
	//日志文件格式
	errorLog = log.New(logfile, "错误是：", log.Ldate|log.Ltime|log.Lshortfile)
		//打开MySQL数据库
		dsn := os.Getenv("MYSQL_DSN")
		if dsn == "" {
			dsn = "root:123456@tcp(127.0.0.1:3306)/library?charset=utf8mb4&parseTime=True&loc=Local"
		}
		db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("无法连接数据库:", err)
	}

	//	用户信息表		书名	读者号	密码	当前借书数	历史借书数	生日	年龄	头像	邮箱	邮箱已验证
	usersTable := `CREATE TABLE IF NOT EXISTS users (
		name VARCHAR(255) NOT NULL DEFAULT '',
		username VARCHAR(255) PRIMARY KEY,
		password VARCHAR(255) NOT NULL DEFAULT '',
		user_cur_lend_amount INT DEFAULT 0,
		user_his_lend_amount INT DEFAULT 0,
		birthday VARCHAR(50) DEFAULT '',
		age INT DEFAULT 0,
		photo VARCHAR(500) DEFAULT '',
		email VARCHAR(255) DEFAULT '',
		email_verified TINYINT(1) DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err = db.Exec(usersTable)
	if err != nil {
		errorLog.Println(err)
	}

	//所有图书表		书名	作者	图书类型	出版社	出版日期	ISBN	封面	简介	价格	数量	可借数量	当前可借数量	是否推荐
	allbooksTable := `CREATE TABLE IF NOT EXISTS all_books (
		title VARCHAR(500) NOT NULL DEFAULT '',
		author VARCHAR(255) DEFAULT '',
		book_type VARCHAR(100) DEFAULT '',
		press VARCHAR(255) DEFAULT '',
		press_date VARCHAR(50) DEFAULT '',
		isbn VARCHAR(100) PRIMARY KEY,
		cover VARCHAR(500) DEFAULT '',
		intro TEXT,
		price DECIMAL(10,2) DEFAULT 0,
		amount INT DEFAULT 0,
		lend_amount INT DEFAULT 0,
		cur_lend_amount INT DEFAULT 0,
		rec_state INT DEFAULT 0
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err = db.Exec(allbooksTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 历史借阅记录表		借阅编号	读者号	书名	isbn	借阅日期	预计归还日期
	lendrecordTable := `CREATE TABLE IF NOT EXISTS lend_records (
		lend_id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL DEFAULT '',
		title VARCHAR(500) NOT NULL DEFAULT '',
		isbn VARCHAR(100) NOT NULL DEFAULT '',
		lend_date VARCHAR(50) DEFAULT '',
		exp_return_date VARCHAR(50) DEFAULT ''
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err = db.Exec(lendrecordTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 归还记录表		归还编号	读者号	书名	isbn	预计归还日期	归还日期	逾期费用
	returnrecordTable := `CREATE TABLE IF NOT EXISTS return_records (
		return_id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL DEFAULT '',
		title VARCHAR(500) NOT NULL DEFAULT '',
		isbn VARCHAR(100) NOT NULL DEFAULT '',
		lend_date VARCHAR(50) DEFAULT '',
		exp_return_date VARCHAR(50) DEFAULT '',
		return_date VARCHAR(50) DEFAULT '',
		late_fee DECIMAL(10,2) DEFAULT 0
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err = db.Exec(returnrecordTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 当前借阅记录表		借阅编号	读者号	书名	isbn	借阅日期	预计归还日期
	curlendrecordTable := `CREATE TABLE IF NOT EXISTS cur_lend_records (
		lend_id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL DEFAULT '',
		title VARCHAR(500) NOT NULL DEFAULT '',
		isbn VARCHAR(100) NOT NULL DEFAULT '',
		lend_date VARCHAR(50) DEFAULT '',
		exp_return_date VARCHAR(50) DEFAULT ''
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err = db.Exec(curlendrecordTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 推荐图书表		isbn	书名	作者	推荐类型	封面	当前可借数量
	recommendbookTable := `CREATE TABLE IF NOT EXISTS recommend_books (
		isbn VARCHAR(100) PRIMARY KEY,
		title VARCHAR(500) NOT NULL DEFAULT '',
		author VARCHAR(255) DEFAULT '',
		rec_type VARCHAR(100) DEFAULT '',
		cover VARCHAR(500) DEFAULT '',
		cur_lend_amount INT DEFAULT 0
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err = db.Exec(recommendbookTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 公告表	公告编号	公告日期	公告标题	公告内容
	noticeTable := `CREATE TABLE IF NOT EXISTS notices (
		notice_id INT AUTO_INCREMENT PRIMARY KEY,
		notice_date VARCHAR(50) DEFAULT '',
		notice_title VARCHAR(500) DEFAULT '',
		notice TEXT
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err = db.Exec(noticeTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 图书调整表		调整编号	调整日期	调整书的书名	调整书的isbn	调整内容
	adjustbookTable := `CREATE TABLE IF NOT EXISTS adjust_books (
		adjust_id INT AUTO_INCREMENT PRIMARY KEY,
		adjust_date VARCHAR(50) DEFAULT '',
		adjust_title VARCHAR(500) DEFAULT '',
		adjust_isbn VARCHAR(100) DEFAULT '',
		adjust_content TEXT
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err = db.Exec(adjustbookTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 用户意见建议表		用户建议编号	用户名	电话	电子邮件	意见建议
	useropinionTable := `CREATE TABLE IF NOT EXISTS user_opinions (
		opinion_id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) DEFAULT '',
		phone VARCHAR(100) DEFAULT '',
		email VARCHAR(255) DEFAULT '',
		idea TEXT
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err = db.Exec(useropinionTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 回复意见建议表	回复意见建议id	回复单位	回复日期	回复内容	回复的用户
	replayopinionsTable := `CREATE TABLE IF NOT EXISTS replay_opinions (
		replay_id INT AUTO_INCREMENT PRIMARY KEY,
		replay_name VARCHAR(255) DEFAULT '',
		replay_date VARCHAR(50) DEFAULT '',
		replay_idea TEXT,
		replay_user VARCHAR(255) DEFAULT ''
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err = db.Exec(replayopinionsTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 汇总表		总书数量	总书当前可借数量	总书待归还数量	总书用户数量
	libraryTable := `CREATE TABLE IF NOT EXISTS library_summary (
		total_books_amount INT DEFAULT 0,
		total_lend_amount INT DEFAULT 0,
		total_return_amount INT DEFAULT 0,
		total_users_amount INT DEFAULT 0
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err = db.Exec(libraryTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 刷新令牌表 (替代旧的 session_state)
	refreshTokenTable := `CREATE TABLE IF NOT EXISTS refresh_tokens (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		token_hash VARCHAR(255) NOT NULL,
		expires_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err = db.Exec(refreshTokenTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 邮箱验证码表
	verificationCodeTable := `CREATE TABLE IF NOT EXISTS verification_codes (
			id INT AUTO_INCREMENT PRIMARY KEY,
			email VARCHAR(255) NOT NULL,
			code VARCHAR(6) NOT NULL,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err = db.Exec(verificationCodeTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 图书内容表（在线阅读用）
	bookContentTable := `CREATE TABLE IF NOT EXISTS book_contents (
			isbn VARCHAR(100) PRIMARY KEY,
			content LONGTEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err = db.Exec(bookContentTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 管理员表
	adminTable := `CREATE TABLE IF NOT EXISTS admin (
			admin_id VARCHAR(100) PRIMARY KEY,
			admin_password VARCHAR(255) NOT NULL DEFAULT '',
			admin_role VARCHAR(50) DEFAULT ''
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err = db.Exec(adminTable)
	if err != nil {
		errorLog.Println(err)
	}
}

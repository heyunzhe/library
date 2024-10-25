package mode

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var errorLog *log.Logger
var db *sql.DB

func Init() {
	var err error
	//创建日志文件
	logfile, err := os.OpenFile("library.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("无法打开日志文件:", err)
	}
	//日志文件格式
	errorLog = log.New(logfile, "错误是：", log.Ldate|log.Ltime|log.Lshortfile)
	//打开数据库
	db, err = sql.Open("sqlite3", "library.db")
	if err != nil {
		log.Fatal(err)
	}
	//	用户信息表		用户名	读者号	密码	当前借书数	历史借书数	生日	年龄	头像
	usersTable := `CREATE TABLE IF NOT EXISTS users (
	name TEXT,
	username TEXT PRIMARY KEY, 
	password TEXT,
	user_cur_lend_amount,
	user_his_lend_amount,
	birthday TEXT,
	age INTEGER,
	photo TEXT
	);`
	_, err = db.Exec(usersTable)
	if err != nil {
		errorLog.Println(err)
	}
	//所有图书表		书名	作者	图书类型	出版社	出版日期	ISBN	封面	简介	价格	数量	可借数量	当前可借数量	是否推荐
	allbooksTable := `CREATE TABLE IF NOT EXISTS all_books (
	title TEXT, 
	author TEXT, 
	book_type TEXT,
	press TEXT,
	press_date TEXT,
	isbn TEXT PRIMARY KEY,
	cover TEXT,
	intro TEXT,
	price INTEGER,
	amount INTEGER,
	lend_amount INTEGER,
	cur_lend_amount INTEGER,
	rec_state INTEGER
	);`
	_, err = db.Exec(allbooksTable)
	if err != nil {
		errorLog.Println(err)
	}
	// 历史借阅记录表		借阅编号	读者号	书名	isbn	借阅日期	预计归还日期
	lendrecordTable := `CREATE TABLE IF NOT EXISTS lend_records (
	lend_id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT,
	title TEXT,
	isbn TEXT,
	lend_date TEXT,
	exp_return_date TEXT
	)`
	_, err = db.Exec(lendrecordTable)
	if err != nil {
		errorLog.Println(err)
	}
	// 归还记录表		归还编号	读者号	书名	isbn	预计归还日期	归还日期	逾期费用
	returnrecordTable := `CREATE TABLE IF NOT EXISTS return_records (
	return_id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT,
	title TEXT,
	isbn TEXT, 
	lend_date TEXT,
	exp_return_date TEXT,
	return_date TEXT,
	late_fee INTEGER
	)`
	_, err = db.Exec(returnrecordTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 当前借阅记录表		借阅编号	读者号	书名	isbn	借阅日期	预计归还日期
	curlendrecordTable := `CREATE TABLE IF NOT EXISTS cur_lend_records (
	lend_id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT,
	title TEXT,
	isbn TEXT,
	lend_date TEXT,
	exp_return_date TEXT
	)`
	_, err = db.Exec(curlendrecordTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 推荐图书表		isbn	书名	作者	推荐类型	封面	当前可借数量
	recommendbookTable := `CREATE TABLE IF NOT EXISTS recommend_books (
	isbn TEXT PRIMARY KEY, 
	title TEXT,
	author TEXT,
	rec_type TEXT,
	cover TEXT,
	cur_lend_amount INTEGER
	)`
	_, err = db.Exec(recommendbookTable)
	if err != nil {
		errorLog.Println(err)
	}
	// 公告表	公告编号	公告日期	公告标题	公告内容
	noticeTable := `CREATE TABLE IF NOT EXISTS notices (
	notice_id INTEGER PRIMARY KEY AUTOINCREMENT,
	notice_date TEXT,
	notice_title TEXT,
	notice TEXT
	)`
	_, err = db.Exec(noticeTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 图书调整表		调整编号	作者	书名	调整类型	调整日期	原因
	// adjustbookTable := `CREATE TABLE IF NOT EXISTS adjust_books (
	// adj_id INTEGER PRIMARY KEY AUTOINCREMENT,
	// author TEXT,
	// title TEXT,
	// adj_type TEXT,
	// adj_date TEXT,
	// cause TEXT
	// )`
	// _, err = db.Exec(adjustbookTable)
	// if err != nil {
	// 	errorLog.Println(adjustbookTable)
	// }

	// 用户意见建议表		用户建议编号	用户名	电话	电子邮件	意见建议
	useropinionTable := `CREATE TABLE IF NOT EXISTS user_opinions (
	opinion_id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	phone TEXT,
	email TEXT,
	idea TEXT
	)`
	_, err = db.Exec(useropinionTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 回复意见建议表	回复意见建议id	回复单位	回复日期	回复内容	回复的用户
	replayopinionsTable := `CREATE TABLE IF NOT EXISTS replay_opinions(
	replay_id INTEGER PRIMARY KEY AUTOINCREMENT, 
	replay_name TEXT,
	replay_date TEXT,
	replay_idea TEXT,
	replay_user TEXT
	)`
	_, err = db.Exec(replayopinionsTable)
	if err != nil {
		errorLog.Println(err)
	}

	// 汇总表		总书数量	总书当前可借数量	总书待归还数量	总书用户数量
	libraryTable := `CREATE TABLE IF NOT EXISTS library_summary (
	total_books_amount INTEGER,
	total_lend_amount INTEGER,
	total_return_amount INTEGER,
	total_users_amount INTEGER
	)`
	_, err = db.Exec(libraryTable)
	if err != nil {
		errorLog.Println(err)
	}

	sessionstate := `CREATE TABLE IF NOT EXISTS session_state (
	session_name TEXT,
	session_id TEXT,
	)`
	_, err = db.Exec(sessionstate)
	if err != nil {
		errorLog.Println(err)
	}

}

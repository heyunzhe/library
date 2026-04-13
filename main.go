package main

import (
	"fmt"
	"librarys/go"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	mode.Init()
	http.HandleFunc("/add/book", mode.AdminAuthMiddleware(mode.AddBooksHandler))          // 添加图书
	http.HandleFunc("/update/book", mode.AdminAuthMiddleware(mode.UpdateBookHandler))     //更新图书
	http.HandleFunc("/view/book", mode.AdminAuthMiddleware(mode.ViewBookHandler))         //查询图书
	http.HandleFunc("/delete/book", mode.AdminAuthMiddleware(mode.DeleteBookHandler))     //删除图书
	http.HandleFunc("/view/user", mode.AdminAuthMiddleware(mode.ViewUserHandler))         //查询用户
	http.HandleFunc("/add/notice", mode.AdminAuthMiddleware(mode.AddNoticeHandler))       //添加公告
	http.HandleFunc("/view/notice", mode.AdminAuthMiddleware(mode.ViewNoticeHandler))     //查询公告
	http.HandleFunc("/update/notice", mode.AdminAuthMiddleware(mode.UpdateNoticeHandler)) //更新公告
	http.HandleFunc("/delete/notice", mode.AdminAuthMiddleware(mode.DeleteNoticeHandler)) //删除公告
	// http.HandleFunc("/sum", mode.SumlibraryHandler)       //查看汇总信息
	http.HandleFunc("/view/useropi", mode.AdminAuthMiddleware(mode.ViewUserOpinionHandler))     //查询用户意见建议
	http.HandleFunc("/replay/useropi", mode.AdminAuthMiddleware(mode.ReplayUserOpinionHandler)) // 回复用户意见建议

	http.HandleFunc("/admin", mode.AdminHandler)                                   //登录后台
	http.HandleFunc("/index", mode.IndexHandler)                                   //进入首页
	http.HandleFunc("/", mode.IndexHandler)                                        //进入首页
	http.HandleFunc("/lend/book", mode.LendBookHandler)                            //进入借书界面
	http.HandleFunc("/about", mode.AboutHandler)                                   //进入关于我们
	http.HandleFunc("/add/useropi", mode.AddUserOpinionHandler)                    //用户上传意见
	http.HandleFunc("/logout", mode.AdminAuthMiddleware(mode.LogoutHandler))       //管理员退出登录
	http.HandleFunc("/login", mode.LoginHandler)                                   //用户登录
	http.HandleFunc("/ulogout", mode.AuthMiddleware(mode.UserLogoutHandler))       //用户退出登录
	http.HandleFunc("/user/library", mode.AuthMiddleware(mode.UserLibraryHandler)) //进入个人中心界面
	http.HandleFunc("/update/user", mode.UpdateUserHandler)                        //更新用户信息
	http.HandleFunc("/reset", mode.ResetpasswordHandler)                           //重置用户密码
	http.HandleFunc("/ranking", mode.RankingHandler)                               //目前用于测试接口
	http.HandleFunc("/return/book", mode.ReturnBookHandler)                        //还书操作

	http.HandleFunc("/lend/records", mode.AdminAuthMiddleware(mode.ViewLendRecords))
	http.HandleFunc("/return/records", mode.AdminAuthMiddleware(mode.ViewReturnRecords))
	http.HandleFunc("/search/book", mode.ViewSearchBookHandler)
	http.HandleFunc("/adjust/book", mode.AdjustBookHandler)
	http.HandleFunc("/view/adjust", mode.ViewAdjustBookHandler)

	http.HandleFunc("/class/search", mode.ClassifySearchHandler)

	fs := http.FileServer(http.Dir("./"))
	http.Handle("/css/", fs)
	http.Handle("/js/", fs)
	http.Handle("/font/", fs)
	http.Handle("/images/", fs)
	http.Handle("/userphoto/", fs)

	fmt.Println("服务器在 http://localhost:8080 上运行")
	fmt.Println("服务器在 http://10.1.10.118:8080 上运行")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

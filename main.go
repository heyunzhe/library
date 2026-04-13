package main

import (
	"fmt"
	"librarys/go"
	"log"
	"net/http"
)

func main() {

	mode.Init()
	mode.InitPrometheus()

	/* ========== metrics ========== */
	http.Handle("/metrics", mode.MetricsHandler())

	/* ========== 管理员接口（带登录 + 监控） ========== */

	http.HandleFunc("/add/book",
		mode.Chain(mode.AddBooksHandler, "AddBook", mode.AdminAuthMiddleware))

	http.HandleFunc("/update/book",
		mode.Chain(mode.UpdateBookHandler, "UpdateBook", mode.AdminAuthMiddleware))

	http.HandleFunc("/view/book",
		mode.Chain(mode.ViewBookHandler, "ViewBook", mode.AdminAuthMiddleware))

	http.HandleFunc("/delete/book",
		mode.Chain(mode.DeleteBookHandler, "DeleteBook", mode.AdminAuthMiddleware))

	http.HandleFunc("/view/user",
		mode.Chain(mode.ViewUserHandler, "ViewUser", mode.AdminAuthMiddleware))

	http.HandleFunc("/add/notice",
		mode.Chain(mode.AddNoticeHandler, "AddNotice", mode.AdminAuthMiddleware))

	http.HandleFunc("/view/notice",
		mode.Chain(mode.ViewNoticeHandler, "ViewNotice", mode.AdminAuthMiddleware))

	http.HandleFunc("/update/notice",
		mode.Chain(mode.UpdateNoticeHandler, "UpdateNotice", mode.AdminAuthMiddleware))

	http.HandleFunc("/delete/notice",
		mode.Chain(mode.DeleteNoticeHandler, "DeleteNotice", mode.AdminAuthMiddleware))

	http.HandleFunc("/view/useropi",
		mode.Chain(mode.ViewUserOpinionHandler, "ViewUserOpinion", mode.AdminAuthMiddleware))

	http.HandleFunc("/replay/useropi",
		mode.Chain(mode.ReplayUserOpinionHandler, "ReplayUserOpinion", mode.AdminAuthMiddleware))

	http.HandleFunc("/admin", mode.AdminHandler)
	http.HandleFunc("/index", mode.IndexHandler)
	http.HandleFunc("/", mode.IndexHandler)

	http.HandleFunc("/lend/book",
		mode.Chain(mode.LendBookHandler, "Lendbook", mode.AdminAuthMiddleware))
	http.HandleFunc("/about", mode.AboutHandler)

	http.HandleFunc("/add/useropi", mode.AddUserOpinionHandler)

	http.HandleFunc("/login", mode.LoginHandler)

	http.HandleFunc("/logout",
		mode.Chain(mode.LogoutHandler, "Logout", mode.AdminAuthMiddleware))

	http.HandleFunc("/ulogout",
		mode.Chain(mode.UserLogoutHandler, "UserLogout", mode.AuthMiddleware))

	http.HandleFunc("/user/library",
		mode.Chain(mode.UserLibraryHandler, "UserLibrary", mode.AuthMiddleware))

	http.HandleFunc("/update/user", mode.UpdateUserHandler)
	http.HandleFunc("/reset", mode.ResetpasswordHandler)
	http.HandleFunc("/ranking", mode.RankingHandler)
	http.HandleFunc("/return/book", mode.ReturnBookHandler)

	http.HandleFunc("/lend/records",
		mode.Chain(mode.ViewLendRecords, "LendRecords", mode.AdminAuthMiddleware))

	http.HandleFunc("/return/records",
		mode.Chain(mode.ViewReturnRecords, "ReturnRecords", mode.AdminAuthMiddleware))

	http.HandleFunc("/search/book",
		mode.Chain(mode.ViewSearchBookHandler, "SearchBook", func(h http.HandlerFunc) http.HandlerFunc {
			return h
		}))

	http.HandleFunc("/adjust/book",
		mode.Chain(mode.AdjustBookHandler, "AdjustBook", func(h http.HandlerFunc) http.HandlerFunc {
			return h
		}))

	http.HandleFunc("/view/adjust",
		mode.Chain(mode.ViewAdjustBookHandler, "ViewAdjust", mode.AdminAuthMiddleware))

	http.HandleFunc("/class/search",
		mode.Chain(mode.ClassifySearchHandler, "ClassSearch", func(h http.HandlerFunc) http.HandlerFunc {
			return h
		}))

	/* ========== 静态资源 ========== */

	fs := http.FileServer(http.Dir("./"))
	http.Handle("/css/", fs)
	http.Handle("/js/", fs)
	http.Handle("/font/", fs)
	http.Handle("/images/", fs)
	http.Handle("/userphoto/", fs)

	fmt.Println("http://localhost:8080")
	fmt.Println("服务在8080端口运行")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

//@title Restful API
//@version 1.0.0
//@description Define an API
//@Schemes http
//@host localhost:8080
//@BasePath /
package main

import (
	"api/controllers"
	"api/driver"
	"api/utils"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

//資料庫
var db *driver.DB

//初始化連線
func init() {
	//read .env file
	gotenv.Load()

	//設定資料庫資訊
	user := driver.MySqlUser{
		Host:     os.Getenv("db_host"), //主機
		MaxIdle:  10,                   //閒置的連接數
		MaxOpen:  10,                   //最大連接數
		User:     os.Getenv("db_user"), //用戶名
		Password: os.Getenv("db_pass"), //密碼
		Database: os.Getenv("db_name"), //資料庫名稱
		Port:     os.Getenv("db_port"), //端口
	}
	//連線
	db = user.Init()
}

func main() {
	//最後必須關閉
	defer db.Close()

	//create router
	//func NewRouter() *Router
	router := mux.NewRouter()
	controller := controllers.Controller{}
	//func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *Route
	//func (r *Router) Methods(methods ...string) *Route
	router.HandleFunc("/signup", controller.Signup(db)).Methods("POST")
	router.HandleFunc("/login", controller.Login(db)).Methods("POST")
	router.HandleFunc("/data", controller.Contacts(db)).Methods("GET")
	router.HandleFunc("/data/{id}", controller.Contact(db)).Methods("GET")

	//func (r *Router) Use(mc MiddlewareChain)
	//attach JWT auth middleware
	router.Use(utils.JwtAuthentication)
	//localhost
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	//connect server
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}

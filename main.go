package main

import (
	"api/controllers"
	"api/driver"
	models "api/model"
	"api/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"src/github.com/subosito/gotenv"
)

var db *driver.DB
var data driver.ColName
var datas []driver.ColName
var controller controllers.Controller

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

	db = user.Init()
}

func main() {
	//create router
	//func NewRouter() *Router
	router := mux.NewRouter()
	//func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *Route
	//func (r *Router) Methods(methods ...string) *Route
	router.HandleFunc("/signup", controller.Signup(db)).Methods("POST")
	router.HandleFunc("/login", controller.Login(db)).Methods("POST")
	router.HandleFunc("/datas", controller.Contacts(db)).Methods("GET")
	router.HandleFunc("/datas/{id}", controller.Contact(db)).Methods("GET")

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

func Contact(w http.ResponseWriter, r *http.Request) {
	var error models.Error

	//return map
	//func Vars(r *http.Request) map[string]string
	params := mux.Vars(r)
	fmt.Println(params)
	id, _ := strconv.Atoi(params["id"])
	fmt.Println(id)
	data, err := db.ReadSome("data", "id", id, datas, data)
	if err != nil {
		error.Message = "Serve error"
		//encode
		utils.SendError(w, http.StatusInternalServerError, error)
		return
	}
	utils.SendSuccess(w, data)
}

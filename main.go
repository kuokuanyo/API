package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"src/github.com/dgrijalva/jwt-go"
)

var db *sql.DB

//user information
type User struct {
	ID       int
	Email    string
	Password string
}

//json-web-token
type JWT struct {
	Token string
}

//error
type Error struct {
	Message string
}

//用戶資料
type MySqlUser struct {
	Host string //主機
	//最大連接數
	MaxIdle  int
	MaxOpen  int
	User     string //用戶名
	Password string //密碼
	Database string //資料庫名稱
	Port     int    //端口
}

//設定資料庫資訊
var user = MySqlUser{
	Host:     "127.0.0.1", //主機
	MaxIdle:  10,          //閒置的連接數
	MaxOpen:  10,          //最大連接數
	User:     "root",      //用戶名
	Password: "asdf4440",  //密碼
	Database: "users",     //資料庫名稱
	Port:     3306,        //端口
}

//初始化連線
func init() {
	var err error

	//完整的資料格式: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	//func Sprintf(format string, a ...interface{}) string
	DataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		user.User,
		user.Password,
		user.Host,
		user.Port,
		user.Database)

	//開啟資料庫連線(sql.Open只是初始化sql.DB物件)
	//func Open(driverName, dataSourceName string) (*DB, error)
	//第一個參數為驅動名稱，第二個參數為資料庫的連結
	db, err = sql.Open("mysql", DataSourceName)
	//檢查錯誤
	if err != nil {
		log.Fatal(err)
	}
	//立即檢查資料庫連線是否可用
	//func (db *DB) Ping() error
	err = db.Ping()
	//檢查錯誤
	if err != nil {
		log.Fatal(err)
	}

	//設定最大連接數
	//SetMaxIdleConns設置閒置的連接數
	db.SetMaxIdleConns(user.MaxIdle)
	//SetMaxOpenConns設置最大打開的連接數，默認值為0代表沒有限制
	db.SetMaxOpenConns(user.MaxOpen)
}

func main() {
	//create router
	//func NewRouter() *Router
	router := mux.NewRouter()
	//func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *Route
	//func (r *Router) Methods(methods ...string) *Route
	router.HandleFunc("/signup", Signup).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")

	//connect server
	log.Fatal(http.ListenAndServe(":8080", router))
}

//response
func SendError(w http.ResponseWriter, status int, error Error) {
	w.WriteHeader(status)
	//encode
	json.NewEncoder(w).Encode(error)
}

//response
func SendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	//encode
	json.NewEncoder(w).Encode(data)
}

//func(http.ResponseWriter, *http.Request)
func Signup(w http.ResponseWriter, r *http.Request) {
	var user User
	var error Error
	//users為資料庫所有會員資料
	users := make([]User, 0)

	//decode(必須是指標)
	//寫入user
	json.NewDecoder(r.Body).Decode(&user)

	//email、password不能為空
	if user.Email == "" {
		error.Message = "E-mail is not empty!"
		SendError(w, http.StatusBadRequest, error)
		return
	}
	if user.Password == "" {
		error.Message = "Password is not empty!"
		SendError(w, http.StatusBadRequest, error)
		return
	}

	//查詢是否使用過資料庫裡的email
	rows, err := db.Query("select * from member;")
	if err != nil {
		log.Fatal(err)
	}
	//使用後必須關閉
	defer rows.Close()
	for rows.Next() {
		u := User{}
		//scan
		err := rows.Scan(&u.ID, &u.Email, &u.Password)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}
	//檢查
	for _, user_db := range users {
		if user.Email == user_db.Email {
			error.Message = "E-mail already taken"
			SendError(w, http.StatusForbidden, error)
			return
		}
	}

	//密碼加密
	//func GenerateFromPassword(password []byte, cost int) ([]byte, error)
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Fatal(err)
	}
	//convert to string
	user.Password = string(hash)

	//將會員資料資料加入資料庫中
	_, err = db.Exec("insert into member (email, password) values(?, ?)",
		user.Email, user.Password)
	if err != nil {
		error.Message = "Server error"
		SendError(w, http.StatusInternalServerError, error)
	}

	//加入資料庫後，密碼改為空白
	user.Password = ""

	SendSuccess(w, "Successfully sign up!")
}

//func(http.ResponseWriter, *http.Request)
func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	var jwt JWT
	var error Error

	//decode
	json.NewDecoder(r.Body).Decode(&user)

	//email、password不能為空
	if user.Email == "" {
		error.Message = "E-mail is not empty!"
		SendError(w, http.StatusBadRequest, error)
		return
	}
	if user.Password == "" {
		error.Message = "Password is not empty!"
		SendError(w, http.StatusBadRequest, error)
		return
	}
	//登入輸入的密碼
	password := user.Password

	//尋找資料庫裡是否有對應的email
	row := db.QueryRow("select * from member where email=?", user.Email)
	//scan
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		//假設找不報此資料
		if err == sql.ErrNoRows {
			error.Message = "The user does not exist!"
			SendError(w, http.StatusBadRequest, error)
			return
		} else {
			log.Fatal(err)
		}
	}
	//資料庫裡的密碼
	hashedpassword := user.Password

	//有可能原始資料庫密碼沒有加密
	if password != hashedpassword {
		//比較密碼是否符合
		//func CompareHashAndPassword(hashedPassword, password []byte) error
		//亂碼的密碼與純文本密碼比較
		err = bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))
		if err != nil {
			error.Message = "Invaild Password"
			SendError(w, http.StatusUnauthorized, error)
			return
		}
	}

	//create token
	token, err := GenerateToken(user)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	jwt.Token = token
	SendSuccess(w, jwt)
}

//json-web-token
func GenerateToken(user User) (string, error) {
	var err error
	s := "kuokuanyo"

	//a jwt
	//header.payload.s
	//func NewWithClaims(method SigningMethod, claims Claims) *Token
	claims := jwt.MapClaims{
		"email":    user.Email,
		"password": user.Password,
		"exp":      time.Now().Add(time.Hour * 1).Unix(), //增添過期時間
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//生成簽名字串(s)
	//func (t *Token) SignedString(key interface{}) (string, error)
	tokenString, err := token.SignedString([]byte(s))
	if err != nil {
		log.Fatal(err)
	}
	return tokenString, nil
}

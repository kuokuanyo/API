package main

import (
	"conn"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"src/github.com/dgrijalva/jwt-go"
	"src/github.com/subosito/gotenv"
)

var db *conn.DB

//user information
type User struct {
	ID       int
	Email    string
	Password string
	Area     string
}

//json-web-token
type JWT struct {
	Token string
}

//error
type Error struct {
	Message string
}

type Ubike struct {
	Name  string
	Area  string
	Total int
}

//初始化連線
func init() {
	gotenv.Load()
	//設定資料庫資訊
	user := conn.MySqlUser{
		Host:     os.Getenv("db_host"), //主機
		MaxIdle:  10,                   //閒置的連接數
		MaxOpen:  10,                   //最大連接數
		User:     os.Getenv("db_user"), //用戶名
		Password: os.Getenv("db_pass"), //密碼
		Database: os.Getenv("db_name"), //資料庫名稱
		Port:     os.Getenv("db_port"), //端口
	}

	connect_db := user.Init()
	db = conn.NewDB(connect_db)
}

func main() {
	//create router
	//func NewRouter() *Router
	router := mux.NewRouter()
	//func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *Route
	//func (r *Router) Methods(methods ...string) *Route
	router.HandleFunc("/signup", Signup).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/datas", Contacts).Methods("GET")
	router.HandleFunc("/datas/{id}", Contacts).Methods("GET")

	//func (r *Router) Use(mc MiddlewareChain)
	//attach JWT auth middleware
	router.Use(JwtAuthentication)

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
	if user.Area == "" {
		error.Message = "Area is not empty!"
		SendError(w, http.StatusBadRequest, error)
		return
	}

	//信箱需有@
	if !strings.Contains(user.Email, "@") {
		error.Message = "Email address is error"
		SendError(w, http.StatusBadRequest, error)
		return
	}
	//密碼長度需大於六字元
	if len(user.Password) < 6 {
		error.Message = "Password length is not enough(6 char)!"
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

func Contacts(w http.ResponseWriter, r *http.Request) {
	var ubikes []Ubike

	//使用資料庫
	db.Use_Db("ubike")
	rows, err := db.Query("select name, area, total from data")
	if err != nil {
		log.Fatal(err)
	}
	//最後必須關閉
	defer rows.Close()
	//處理每一行
	for rows.Next() {
		var ubike Ubike
		err := rows.Scan(&ubike.Name, &ubike.Area, &ubike.Total)
		if err != nil {
			log.Fatal(err)
		}
		ubikes = append(ubikes, ubike)
	}
	SendSuccess(w, ubikes)
}

//json-web-token
func GenerateToken(user User) (string, error) {
	var err error
	s := os.Getenv("token_password")

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

//jwt驗證
func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorObject Error
		//從header中獲取token
		authHeader := r.Header.Get("Authorization")
		//不需要驗證的路徑
		paths := []string{"/signup", "/login"}
		//current request path
		requestPath := r.URL.Path

		//不須驗證，直接執行
		for _, path := range paths {
			if path == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		//if authHeader is empty
		if authHeader == "" {
			errorObject.Message = "Missing auth token!"
			SendError(w, http.StatusForbidden, errorObject)
			return
		}

		//split string
		splitted := strings.Split(authHeader, " ")

		//if length is not 2
		if len(splitted) != 2 {
			errorObject.Message = "Invaild token."
			SendError(w, http.StatusUnauthorized, errorObject)
			return
		}

		//取第二個位置的值(token value)
		authHeader = splitted[1]

		//jwt解析並驗證
		//func Parse(tokenString string, keyFunc Keyfunc) (*Token, error)
		//type Keyfunc func(*Token) (interface{}, error)
		/*
			type Token struct {
			Raw       string                 // The raw token.  Populated when you Parse a token
			Method    SigningMethod          // The signing method used or to be used
			Header    map[string]interface{} // The first segment of the token
			Claims    Claims                 // The second segment of the token
			Signature string                 // The third segment of the token.  Populated when you Parse a token
			Valid     bool                   // Is the token valid?  Populated when you Parse/Verify a token
			}
		*/
		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error.")
			}
			return []byte(os.Getenv("token_password")), nil
		})
		if err != nil {
			errorObject.Message = err.Error()
			SendError(w, http.StatusUnauthorized, errorObject)
			return
		}

		//if token is vaild, return true
		if token.Valid {
			//通過驗證
			next.ServeHTTP(w, r)
		} else {
			errorObject.Message = err.Error()
			SendError(w, http.StatusUnauthorized, errorObject)
			return
		}
	})
}

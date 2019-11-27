package driver

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

//用戶資料
type MySqlUser struct {
	Host string //主機
	//最大連接數
	MaxIdle  int
	MaxOpen  int
	User     string //用戶名
	Password string //密碼
	Database string //資料庫名稱
	Port     string //端口
}

type DB struct {
	*sql.DB
}

/*
func NewDB(d *sql.DB) *DB {
	return &DB{d}
}
*/

//定義資料庫連線連線
//完整的資料格式: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
//mehtod
func (msu *MySqlUser) Init() *DB {

	//資料庫連結字串
	//func Sprintf(format string, a ...interface{}) string
	DataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		msu.User,
		msu.Password,
		msu.Host,
		msu.Port,
		msu.Database)

	//開啟資料庫連線(sql.Open只是初始化sql.DB物件)
	//func Open(driverName, dataSourceName string) (*DB, error)
	//第一個參數為驅動名稱，第二個參數為資料庫的連結
	db, err := sql.Open("mysql", DataSourceName)
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
	db.SetMaxIdleConns(msu.MaxIdle)
	//SetMaxOpenConns設置最大打開的連接數，默認值為0代表沒有限制
	db.SetMaxOpenConns(msu.MaxOpen)
	//無錯誤
	return &DB{db}
}

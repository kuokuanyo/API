package main

import (
	"conn"
	"fmt"
)

//設定資料庫資訊
var user = conn.MySqlUser{
	Host:     "127.0.0.1", //主機
	MaxIdle:  10,          //閒置的連接數
	MaxOpen:  10,          //最大連接數
	User:     "root",      //用戶名
	Password: "asdf4440",  //密碼
	Database: "ubike",     //資料庫名稱
	Port:     3306,        //端口
}

/*
//建立查詢欄位
var (
	名稱	int
	名稱	string
	名稱	bool
)

//上面查詢欄位名稱等於此[]string{}的變數名稱
//須為字串
var s = []string{上列設定名稱}
*/
func main() {

	//建立初始化連線
	connect_db := user.Init()

	db := conn.NewDB(connect_db)

	//最後必須關閉
	defer db.Close()

	//建立資料庫
	db.CreateDb("users")

	//使用資料庫
	db.Use_Db("users")

	//建立資料表
	stmt, err := db.Prepare("CREATE Table member(id int NOT NULL AUTO_INCREMENT, email varchar(50), password varchar(100), PRIMARY KEY (id));")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}

	//插入數值
	db.Insert("member", "email", "test1@example.com", "password", "12345")
	db.Insert("member", "email", "test2@example.com", "password", "abcd")
	db.Insert("member", "email", "test3@example.com", "password", "12345")
	db.Insert("member", "email", "test4@example.com", "password", "abcd")
	db.Insert("member", "email", "test5@example.com", "password", "12345")
	db.Insert("member", "email", "test6@example.com", "password", "abcd")
	db.Insert("member", "email", "test7@example.com", "password", "12345")
	db.Insert("member", "email", "test8@example.com", "password", "abcd")
	db.Insert("member", "email", "test9@example.com", "password", "12345")
	db.Insert("member", "email", "test10@example.com", "password", "abcd")

	/*
		//更改數值
		conn.Update_db(db, 資料庫名稱, 設定欄位名稱, 設定新數值, 更改的欄位, 更改欄位的數值)

		//刪除資料庫
		conn.Delete_Db(db , 資料庫名稱)

		//刪除資料表
		conn.Delete_Tb(db, 資料庫名稱)

		//讀取資料
		//第三個與後面參數長度必須相同
		conn.Read(db, 資料庫名稱, s, 設定的變數(var))
	*/
}

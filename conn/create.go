package conn

import (
	"fmt"
	"log"
)

//新建資料庫
func (db DB) CreateDb(DbName string) {

	//字串(檢查資料庫是否已經存在)
	//func Sprintf(format string, a ...interface{}) string
	CheckDb := fmt.Sprintf("DROP DATABASE IF exists %s", DbName)

	//func (db *DB) Exec(query string, args ...interface{}) (Result, error)
	_, err := db.Exec(CheckDb)
	//檢查錯誤
	if err != nil {
		log.Fatal(err)
	}

	//字串(建立資料庫)
	//func Sprintf(format string, a ...interface{}) string
	CreateName := fmt.Sprintf("CREATE DATABASE %s", DbName)

	//func (db *DB) Exec(query string, args ...interface{}) (Result, error)
	_, err = db.Exec(CreateName)
	//檢查錯誤
	if err != nil {
		log.Fatal(err)
	}
}

//使用資料庫
func (db DB) Use_Db(DbName string) {

	//字串(使用資料庫)
	UseName := fmt.Sprintf("USE %s;", DbName)

	//func (db *DB) Exec(query string, args ...interface{}) (Result, error)
	_, err := db.Exec(UseName)
	//檢查錯誤
	if err != nil {
		log.Fatal(err)
	}
}

//建立資料表
//args索引奇數為欄位名稱，偶數為欄位類型
//args長度必須為偶數(欄位名稱與類型為一組))
func (db DB) CreateTable(TableName string, args ...string) {

	//檢查是否有存在的table
	//func Sprintf(format string, a ...interface{}) string
	CheckTable := fmt.Sprintf("DROP TABLE IF exists %s", TableName)
	//如果有，刪除
	//func (db *DB) Exec(query string, args ...interface{}) (Result, error)
	_, err := db.Exec(CheckTable)
	//檢查錯誤
	if err != nil {
		log.Fatal(err)
	}

	//欄位長度
	n := len(args)
	//args長度必須為偶數
	if n%2 != 0 {
		fmt.Println("Number of Parameters is wrong.")
	}

	//字串(建立資料表)
	//CREATE TABLE tablename(col_name col_type);
	CreateDb_str := ""

	//加入欄位及類型
	//奇數為欄位名稱
	//偶數為欄位類型
	for i := 0; i < n; i++ {
		//開頭
		if i == 0 {
			CreateDb_str += fmt.Sprintf("CREATE TABLE %s(%s ", TableName, args[i])
		} else if i == n-1 {
			CreateDb_str += fmt.Sprintf("%s)", args[i])
		} else if (i+1)%2 == 0 { //偶數
			CreateDb_str += fmt.Sprintf("%s, ", args[i])
		} else if (i+1)%2 == 1 { //奇數
			CreateDb_str += fmt.Sprintf("%s ", args[i])
		}
	}

	//建立資料表
	//func (db *DB) Exec(query string, args ...interface{}) (Result, error)
	_, err = db.Exec(CreateDb_str)
	//檢查錯誤
	if err != nil {
		log.Fatal(err)
	}
}

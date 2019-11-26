//插入數值
package conn

import (
	"fmt"
	"log"
)

//插入值
//INSERT INTO tablename(col_name...) values(?...)
//args為插入的數值
func (db DB) Insert(TableName string, args ...interface{}) {

	//加入欄位的長度
	n := len(args)

	//檢查欄位數量是否與加入值數量是否符合
	//args參數數量為偶數個(欄位及插入設值為一組)
	if n%2 != 0 {
		fmt.Println("Number of Parameters is wrong.")
	}

	//字串(插入數值)
	//INSERT INTO tablename(col_name...) values(?...)
	Insert_str := ""

	//args的奇數位置為欄位名稱(go中的偶數位置)))
	//加入欄位
	for i := 0; i < n; i++ {
		if n == 2 {
			Insert_str += fmt.Sprintf("INSERT INTO %s(%s) values(", TableName, args[i])
			break
		} else if i == 0 { //第一個數值
			Insert_str += fmt.Sprintf("INSERT INTO %s(%s, ", TableName, args[i])
		} else if i == n-2 { //最後一個數值
			Insert_str += fmt.Sprintf("%s) values(", args[i])
		} else if i%2 == 0 {
			Insert_str += fmt.Sprintf("%s, ", args[i])
		}
	}

	//新增占位符
	for i := 0; i < n/2; i++ {
		if n == 1 {
			Insert_str += "?)"
		} else if i == (n/2)-1 {
			Insert_str += "?)"
		} else {
			Insert_str += "?, "
		}
	}


	/*
		使用預編譯語句(Prepared Statement)
		可實現自定義參數查詢
		可防止mysql被攻擊
		比手動拼接更有效率
	*/
	//func (db *DB) Prepare(query string) (*Stmt, error)
	//Prepaer method 是為了之後的運行先準備好的語法
	//如需要運行此語法，執行stmt.Exec()
	//可以同時運行多個查詢

	stmt, err := db.Prepare(Insert_str)
	defer stmt.Close() //一定要關閉
	//檢查錯誤
	if err != nil {
		log.Fatal(err)
	}

	//因插入的型態不同，因此建立[]interface{}
	//放入插入數值
	var t []interface{}

	//args偶數索引為插入數值(go的索引從0開始)
	for i := 0; i < n; i++ {
		if (i+1)%2 == 0 {
			t = append(t, args[i])
		}
	}

	//插入
	_, err = stmt.Exec(t...)
	//檢查錯誤
	if err != nil {
		log.Fatal(err)
	}

}

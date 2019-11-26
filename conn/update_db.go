//更新資料
package conn

import (
	"fmt"
	"log"
)

//函式(更新資料表)
//args奇數索引為欄位，偶數為更改數值
//args總共為四個參數，前兩個為設定數值，後兩個為更改條件
//update tablename set ... where ...
func (db DB) Update_db(TableName string, args ...interface{}) {

	//args數量為四
	if len(args) != 4 {
		fmt.Println("Parameters length are wrong.")
	}

	//更新字串
	//使用佔位符
	Update_str := fmt.Sprintf(`UPDATE %s set %s = ? WHERE %s = ?`,
		TableName, args[0], args[2])

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
	stmt, err := db.Prepare(Update_str)
	defer stmt.Close() //一定要關閉
	//檢查錯誤
	if err != nil {
		log.Fatal(err)
	}

	//執行
	_, err = stmt.Exec(args[1], args[3])
	//檢查錯誤
	if err != nil {
		log.Fatal(err)
	}
}

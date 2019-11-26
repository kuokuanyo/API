//讀取數據
package conn

import (
	"fmt"
	"log"
)

//function read data
//SELECT col_name FROM tablename;
//[]string為查詢欄位
//args欄位變數類型
//s 與 args 長度需相同
func (db DB) Read(TableName string, s []string, args ...interface{}) {

	//長度需相同
	//args長度
	if len(s) != len(args) {
		fmt.Println("Parameters length are wrong.")
	}
	n := len(s)

	//讀取數據字串
	//"SELECT col from tablename;"
	Read_str := ""
	//添加字串
	for i := 0; i < n; i++ {
		if n == 1 {
			Read_str += fmt.Sprintf("SELECT %s FROM %s;", s[i], TableName)
			break
		} else if i == 0 {
			Read_str += fmt.Sprintf("SELECT %s, ", s[i])
		} else if i == n-1 {
			Read_str += fmt.Sprintf("%s FROM %s;", s[i], TableName)
		} else {
			Read_str += fmt.Sprintf("%s, ", s[i])
		}
	}

	//讀取
	//查詢多條
	//func (db *DB) Query(query string, args ...interface{}) (*Rows, error)
	rows, err := db.Query(Read_str)
	//defer 關閉查詢
	//一定要關閉(延遲)
	defer rows.Close()
	//檢查錯誤
	if err != nil {
		log.Fatal(err)
	}

	//建立指標的[]interface{}
	//rows.Scan函式需使用
	var p []interface{}
	for i := 0; i < n; i++ {
		p = append(p, &args[i])
	}

	//格式化字串
	format_str := "Data: "
	for i := 0; i < n; i++ {
		if n == 1 {
			format_str += "%s\n"
			break
		} else if i == n-1 {
			format_str += "%s\n"
		} else {
			format_str += "%s, "
		}
	}

	//處理每一行
	//Next method 迭代查詢資料，回傳bool
	//func (rs *Rows) Next() bool
	for rows.Next() {

		//Scan method方法用來讀取每一列的值
		//func (rs *Rows) Scan(dest ...interface{}) error
		if err := rows.Scan(p...); err != nil {
			log.Fatal(err)
		}

		//"%s, %s...\n"印出data
		fmt.Printf(format_str, args...)
	}

	//在迴圈中是否有錯誤
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

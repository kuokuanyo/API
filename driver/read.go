//讀取數據
package driver

import (
	models "api/model"
	"fmt"
)

var user models.User

//查詢欄位名稱
type ColName struct {
	Id   int
	Math int
	Eng  int
}

//function read all data
//SELECT col_name FROM tablename;
//args為scan的欄位名稱
func (db DB) ReadAll(TableName string, datas []ColName, data ColName) ([]ColName, error) {

	//讀取數據字串
	//"SELECT col from tablename;"
	Read_str := fmt.Sprintf("SELECT * FROM %s", TableName)

	//讀取
	//查詢多條
	//func (db *DB) Query(query string, args ...interface{}) (*Rows, error)
	rows, err := db.Query(Read_str)
	//defer 關閉查詢
	//一定要關閉(延遲)
	defer rows.Close()
	//檢查錯誤
	if err != nil {
		return datas, err
	}

	//處理每一行
	//Next method 迭代查詢資料，回傳bool
	//func (rs *Rows) Next() bool
	for rows.Next() {
		//Scan method方法用來讀取每一列的值
		//func (rs *Rows) Scan(dest ...interface{}) error
		if err := rows.Scan(&data.Id, &data.Math, &data.Eng); err != nil {
			return datas, err
		}
		datas = append(datas, data)
	}
	//在迴圈中是否有錯誤
	if err := rows.Err(); err != nil {
		return datas, err
	}
	return datas, nil
}

//function read all data
//SELECT col_name FROM tablename;
//args為scan的欄位名稱
func (db DB) ReadSome(TableName string, col string, value int, datas []ColName, data ColName) ([]ColName, error) {

	//讀取數據字串
	//"SELECT col from tablename;"
	Read_str := fmt.Sprintf("SELECT * FROM %s WHERE %s=%d", TableName, col, value)

	//讀取
	//查詢多條
	//func (db *DB) Query(query string, args ...interface{}) (*Rows, error)
	rows, err := db.Query(Read_str)
	//defer 關閉查詢
	//一定要關閉(延遲)
	//defer rows.Close()
	//檢查錯誤
	if err != nil {
		return datas, err
	}

	//處理每一行
	//Next method 迭代查詢資料，回傳bool
	//func (rs *Rows) Next() bool
	for rows.Next() {
		//Scan method方法用來讀取每一列的值
		//func (rs *Rows) Scan(dest ...interface{}) error
		if err := rows.Scan(&data.Id, &data.Math, &data.Eng); err != nil {
			return datas, err
		}
		datas = append(datas, data)
	}

	//在迴圈中是否有錯誤
	if err := rows.Err(); err != nil {
		return datas, err
	}
	return datas, nil
}

//查詢單一
func (db DB) ReadOne(TableName string, user models.User, col string) (models.User, error) {
	//讀取數據字串
	//"SELECT col from tablename where ;"
	Read_str := fmt.Sprintf("SELECT * FROM %s where %s=?;", TableName, col)

	//讀取
	//func (db *DB) Query(query string, args ...interface{}) (*Rows, error)
	row := db.QueryRow(Read_str, user.Email)
	//defer 關閉查詢
	//一定要關閉(延遲)
	if err := row.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		return user, err
	}
	return user, nil
}

//刪除資料庫表
package driver

import (
	"fmt"
	"log"
)

//funciton(delete database)
func (db DB) Delete_Db(DbName string) {

	//刪除字串
	Delete := fmt.Sprintf("DROP DATABASE %s", DbName)

	//刪除
	db.Exec(Delete)
}

//function(delete table)
func (db DB) Delete_Tb(TableName string) {

	//刪除字串
	Delete := fmt.Sprintf("DROP TABLE %s", TableName)

	//刪除
	db.Exec(Delete)
}

//function(delete table)
func (db DB) DeleteValue(TableName string, colname string, value int) error {

	//刪除字串
	Delete := fmt.Sprintf("DELETE FROM %s where %s=%d", TableName, colname, value)

	//刪除
	_, err := db.Exec(Delete)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

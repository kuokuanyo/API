//刪除資料庫表
package conn

import (
	"fmt"
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

package userRepository

import "api/driver"

func (u UserRepository) Contacts(db *driver.DB, datas []driver.ColName, data driver.ColName) ([]driver.ColName, error) {
	datas, err := db.ReadAllData("data", datas, data)
	return datas, err
}

func (u UserRepository) Contact(db *driver.DB, datas []driver.ColName, data driver.ColName, value string) ([]driver.ColName, error) {
	datas, err := db.ReadSomeData("data", "id", value, datas, data)
	return datas, err
}

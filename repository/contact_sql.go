package userRepository

import "api/driver"

func (u UserRepository) Contacts(db *driver.DB, datas []driver.ColName, data driver.ColName) ([]driver.ColName, error) {
	datas, err := db.ReadAll("data", datas, data)
	return datas, err
}

/*
func (u UserRepository) Contact(db *driver.DB, datas []driver.ColName, data driver.ColName, id string) ([]driver.ColName, error) {
	datas, err := db.ReadSome("data", "id", id, datas, data)
	return datas, err
}
*/

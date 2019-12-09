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

func (u UserRepository) Insert(db *driver.DB, data driver.ColName) error {
	err := db.Insert("data", "id", data.Id, "math", data.Math, "eng", data.Eng)
	return err
}

func (u UserRepository) UpdateMath(db *driver.DB, data driver.ColName) error {
	err := db.Update_db("data", "math", data.Math, "id", data.Id)
	return err
}

func (u UserRepository) UpdateEng(db *driver.DB, data driver.ColName) error {
	err := db.Update_db("data", "eng", data.Eng, "id", data.Id)
	return err
}

func (u UserRepository) DeleteID(db *driver.DB, data driver.ColName) error {
	err := db.DeleteValue("data", "id", data.Id)
	return err
}

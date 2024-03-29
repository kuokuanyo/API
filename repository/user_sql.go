package userRepository

import (
	"api/driver"
	models "api/model"
)

type UserRepository struct{}

func (u UserRepository) CheckSignup(db *driver.DB, users []models.User, user models.User) ([]models.User, error) {
	users, err := db.ReadAllUser("members", users, user)
	return users, err
}

func (u UserRepository) InsertSignup(db *driver.DB, user models.User) error {
	//將會員資料資料加入資料庫中
	//奇數參數為欄位名稱
	//偶數為插入值
	err := db.Insert("members", "email", user.Email, "password", user.Password)
	return err
}

func (u UserRepository) Login(db *driver.DB, user models.User) (models.User, error) {
	user, err := db.ReadOneUser("members", user, "email")
	return user, err
}

func (u UserRepository) AdminLogin(db *driver.DB, user models.User) (models.User, error) {
	user, err := db.ReadOneUser("admin", user, "email")
	return user, err
}

package userRepository

import (
	"api/driver"
	models "api/model"
)

type UserRepository struct{}

func (u UserRepository) CheckSignup(db *driver.DB, users []models.User, user models.User) ([]models.User, error) {
	users, err := db.ReadAllUser("member", users, user)
	return users, err
}

func (u UserRepository) InsertSignup(db *driver.DB, user models.User) error {
	//將會員資料資料加入資料庫中
	//奇數參數為欄位名稱
	//偶數為插入值
	err := db.Insert("member", "email", user.Email, "password", user.Password)
	return err
}

func (u UserRepository) Login(db *driver.DB, user models.User) (models.User, error) {
	user, err := db.ReadOneUser("member", user, "email")
	return user, err
}

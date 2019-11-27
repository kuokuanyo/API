package userRepository

import (
	"api/driver"
	models "api/model"
)

type UserRepository struct{}

func (u UserRepository) CheckSignup(db *driver.DB, user models.User) ([]models.User, error) {
	var users []models.User
	//查詢是否使用過資料庫裡的email
	rows, err := db.Query("select * from member;")
	if err != nil {
		return users, err
	}
	//使用後必須關閉
	defer rows.Close()
	for rows.Next() {
		//scan
		err := rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u UserRepository) InsertSignup(db *driver.DB, user models.User) (models.User, error) {
	//將會員資料資料加入資料庫中
	_, err := db.Exec("insert into member (email, password) values(?, ?)",
		user.Email, user.Password)
	return user, err
}

func (u UserRepository) Login(db *driver.DB, user models.User) (models.User, error) {
	user, err := db.ReadOne("member", user, "email")
	return user, err
}

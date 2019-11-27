package controllers

import (
	"api/driver"
	models "api/model"
	userRepository "api/repository"
	"api/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Controller struct{}

//func(http.ResponseWriter, *http.Request)
func (c Controller) Signup(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error
		//users為資料庫所有會員資料
		users := make([]models.User, 0)

		//decode(必須是指標)
		//寫入user
		json.NewDecoder(r.Body).Decode(&user)

		//email、password不能為空
		if user.Email == "" {
			error.Message = "E-mail is not empty!"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}
		if user.Password == "" {
			error.Message = "Password is not empty!"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		//信箱需有@
		if !strings.Contains(user.Email, "@") {
			error.Message = "Email address is error"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}
		//密碼長度需大於六字元
		if len(user.Password) < 6 {
			error.Message = "Password length is not enough(6 char)!"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		userRepo := userRepository.UserRepository{}
		users, err := userRepo.CheckSignup(db, user)
		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
		}
		//檢查
		for _, user_db := range users {
			if user.Email == user_db.Email {
				error.Message = "E-mail already taken"
				utils.SendError(w, http.StatusForbidden, error)
				return
			}
		}

		//密碼加密
		//func GenerateFromPassword(password []byte, cost int) ([]byte, error)
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			log.Fatal(err)
		}
		//convert to string
		user.Password = string(hash)

		user, err = userRepo.InsertSignup(db, user)
		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		//加入資料庫後，密碼改為空白
		user.Password = ""

		utils.SendSuccess(w, "Successfully sign up!")
	}
}

//func(http.ResponseWriter, *http.Request)
func (c Controller) Login(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var err error
		var jwt models.JWT
		var error models.Error

		//decode
		json.NewDecoder(r.Body).Decode(&user)

		//email、password不能為空
		if user.Email == "" {
			error.Message = "E-mail is not empty!"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}
		if user.Password == "" {
			error.Message = "Password is not empty!"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}
		//登入輸入的密碼
		password := user.Password

		userRepo := userRepository.UserRepository{}
		user, err = userRepo.Login(db, user)
		if err != nil {
			//假設找不報此資料
			if err == sql.ErrNoRows {
				error.Message = "The user does not exist!"
				utils.SendError(w, http.StatusBadRequest, error)
				return
			} else {
				error.Message = "Server error"
				utils.SendError(w, http.StatusInternalServerError, error)
				return
			}
		}
		//資料庫裡的密碼
		hashedpassword := user.Password

		//有可能原始資料庫密碼沒有加密
		if password != hashedpassword {
			//比較密碼是否符合
			//func CompareHashAndPassword(hashedPassword, password []byte) error
			//亂碼的密碼與純文本密碼比較
			err = bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))
			if err != nil {
				error.Message = "Invaild Password"
				utils.SendError(w, http.StatusUnauthorized, error)
				return
			}
		}

		//create token
		token, err := utils.GenerateToken(user)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		jwt.Token = token
		utils.SendSuccess(w, jwt)
	}
}

package controllers

import (
	"api/driver"
	models "api/model"
	userRepository "api/repository"
	"api/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

//@Summary get all datas
//@Tags Data
//@Description 取得所有資料
//@Accept json
//@Produce json
//@Success 200 {object} driver.ColName "get all datas"
//@Failure 500 {object} models.Error "Serve(database) error"
//@Router /data [get]
func (c Controller) Contacts(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data driver.ColName
		var datas []driver.ColName
		var error models.Error

		userRepo := userRepository.UserRepository{}
		datas, err := userRepo.Contacts(db, datas, data)
		if err != nil {
			error.Message = "Serve(database) error"
			//encode
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}
		utils.SendSuccess(w, datas)
	}
}

// @Summary get some datas
// @Tags Data
// @Description 取得部分資料
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} driver.ColName "get some datas(from database)"
// @Failure 500 {object} models.Error "Serve(database) error"
// @Router /data/{id} [get]
func (c Controller) Contact(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data driver.ColName
		var datas []driver.ColName
		var error models.Error

		//return map
		//func Vars(r *http.Request) map[string]string
		params := mux.Vars(r)
		//id, _ := strconv.Atoi(params["id"])

		userRepo := userRepository.UserRepository{}
		datas, err := userRepo.Contact(db, datas, data, params["id"])

		if err != nil {
			error.Message = "Serve error"
			//encode
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}
		utils.SendSuccess(w, datas)
	}
}

// @Summary add datas
// @Tags Database(admin)
// @Description 管理員增加資料
// @Accept json
// @Produce json
// @Param information body model.user true "個人資料"
// @Param id path int true "ID"
// @Param math path int true "math"
// @Param eng path int true "eng"
// @Success 200 {object} driver.ColName "data"
// @Failure 400 {object} models.Error "email or password error"
// @Failure 401 {object} models.Error "Invaild Password"
// @Failure 500 {object} models.Error "Serve(database) error"
// @Router /db/add/{id}/{math}/{eng} [post]
func (c Controller) Insert(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var data driver.ColName
		var error models.Error

		params := mux.Vars(r)
		//印出url的數值
		data.Id, _ = strconv.Atoi(params["id"])
		data.Math, _ = strconv.Atoi(params["math"])
		data.Eng, _ = strconv.Atoi(params["eng"])

		//decode
		json.NewDecoder(r.Body).Decode(&user)
		//信箱密碼驗證
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
		user, err := userRepo.AdminLogin(db, user)
		if err != nil {
			//假設找不報此資料
			if err == sql.ErrNoRows {
				error.Message = "The user does not exist!"
				utils.SendError(w, http.StatusBadRequest, error)
				return
			} else {
				error.Message = "Server(database) error"
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

		userRepo = userRepository.UserRepository{}
		err = userRepo.Insert(db, data)

		if err != nil {
			error.Message = "Serve error"
			//encode
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}
		utils.SendSuccess(w, data)
	}
}

// @Summary update data
// @Tags Database(admin)
// @Description 修改數學成績
// @Accept json
// @Produce json
// @Param information body model.user true "個人資料"
// @Param id path int true "ID"
// @Param math path int true "math"
// @Success 200 {string} json "Successful"
// @Failure 400 {object} models.Error "email or password error"
// @Failure 401 {object} models.Error "Invaild Password"
// @Failure 500 {object} models.Error "Serve(database) error"
// @Router /db/update/{id}/{math} [put]
func (c Controller) UpdateMath(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var data driver.ColName
		var error models.Error

		params := mux.Vars(r)
		//印出url的數值
		data.Id, _ = strconv.Atoi(params["id"])
		data.Math, _ = strconv.Atoi(params["math"])

		//decode
		json.NewDecoder(r.Body).Decode(&user)
		//信箱密碼驗證
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
		user, err := userRepo.AdminLogin(db, user)
		if err != nil {
			//假設找不報此資料
			if err == sql.ErrNoRows {
				error.Message = "The user does not exist!"
				utils.SendError(w, http.StatusBadRequest, error)
				return
			} else {
				error.Message = "Server(database) error"
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

		userRepo = userRepository.UserRepository{}
		err = userRepo.UpdateMath(db, data)

		if err != nil {
			error.Message = "Serve error"
			//encode
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}
		utils.SendSuccess(w, "Successful")
	}
}

// @Summary update data
// @Tags Database(admin)
// @Description 修改英文成績
// @Accept json
// @Produce json
// @Param information body model.user true "個人資料"
// @Param id path int true "ID"
// @Param eng path int true "eng"
// @Success 200 {string} json "Successful"
// @Failure 400 {object} models.Error "email or password error"
// @Failure 401 {object} models.Error "Invaild Password"
// @Failure 500 {object} models.Error "Serve(database) error"
// @Router /db/update/{id}/{eng} [put]
func (c Controller) UpdateEng(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var data driver.ColName
		var error models.Error

		params := mux.Vars(r)
		//印出url的數值
		data.Id, _ = strconv.Atoi(params["id"])
		data.Math, _ = strconv.Atoi(params["eng"])

		//decode
		json.NewDecoder(r.Body).Decode(&user)
		//信箱密碼驗證
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
		user, err := userRepo.AdminLogin(db, user)
		if err != nil {
			//假設找不報此資料
			if err == sql.ErrNoRows {
				error.Message = "The user does not exist!"
				utils.SendError(w, http.StatusBadRequest, error)
				return
			} else {
				error.Message = "Server(database) error"
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

		userRepo = userRepository.UserRepository{}
		err = userRepo.UpdateMath(db, data)

		if err != nil {
			error.Message = "Serve error"
			//encode
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}
		utils.SendSuccess(w, "Successful")
	}
}

// @Summary delete data
// @Tags Database(admin)
// @Description 刪除資料
// @Accept json
// @Produce json
// @Param information body model.user true "個人資料"
// @Param id path int true "ID"
// @Success 200 {string} json "Successful"
// @Failure 400 {object} models.Error "email or password error"
// @Failure 401 {object} models.Error "Invaild Password"
// @Failure 500 {object} models.Error "Serve(database) error"
// @Router /db/delete/{id} [delete]
func (c Controller) DeleteID(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var data driver.ColName
		var error models.Error

		params := mux.Vars(r)
		//印出url的數值
		data.Id, _ = strconv.Atoi(params["id"])

		//decode
		json.NewDecoder(r.Body).Decode(&user)
		//信箱密碼驗證
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
		user, err := userRepo.AdminLogin(db, user)
		if err != nil {
			//假設找不報此資料
			if err == sql.ErrNoRows {
				error.Message = "The user does not exist!"
				utils.SendError(w, http.StatusBadRequest, error)
				return
			} else {
				error.Message = "Server(database) error"
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

		userRepo = userRepository.UserRepository{}
		err = userRepo.DeleteID(db, data)

		if err != nil {
			error.Message = "Serve error"
			//encode
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}
		utils.SendSuccess(w, "Successful")
	}
}

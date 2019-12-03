package controllers

import (
	"api/driver"
	models "api/model"
	userRepository "api/repository"
	"api/utils"
	"net/http"

	"github.com/gorilla/mux"
)

//@Summary get all datas
//@Tags Data
//@Description 取得所有資料
//@Accept json
//@Produce json
//@Param token body string true "安全碼"
//@Success 200 {object} driver.ColName "get all datas"
//@Failure 500 {object} models.Error "Serve(database) error"
//@Router /contacts [get]
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
// @Param token body string true "安全碼"
// @Success 200 {object} driver.ColName "get some datas(from database)"
// @Failure 500 {object} models.Error "Serve(database) error"
// @Router /contacts/{id} [get]
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

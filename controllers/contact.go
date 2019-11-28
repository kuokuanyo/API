package controllers

import (
	"api/driver"
	models "api/model"
	userRepository "api/repository"
	"api/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func (c Controller) Contacts(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data driver.ColName
		var datas []driver.ColName
		var error models.Error

		userRepo := userRepository.UserRepository{}
		datas, err := userRepo.Contacts(db, datas, data)
		if err != nil {
			error.Message = "Serve error"
			//encode
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}
		utils.SendSuccess(w, datas)
	}
}

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

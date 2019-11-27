package controllers

import (
	"api/driver"
	models "api/model"
	userRepository "api/repository"
	"api/utils"
	"fmt"
	"net/http"
	"strconv"

	"src/github.com/gorilla/mux"
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
		fmt.Println(params)
		id, _ := strconv.Atoi(params["id"])

		datas, err := db.ReadSome("data", "id", id, datas, data)
		fmt.Println(data)
		fmt.Println(err)
		if err != nil {
			error.Message = "Serve error"
			//encode
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}
		utils.SendSuccess(w, data)
	}
}

/*
func (c Controller) Contact(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data driver.ColName
		var datas []driver.ColName
		var error models.Error
		var err error
		//return map
		//func Vars(r *http.Request) map[string]string
		params := mux.Vars(r)
		//id, _ := strconv.Atoi(params["area"])

		datas, err = db.ReadSome("data", "id", params["id"], datas, data)
		if err != nil {
			error.Message = "Serve error"
			//encode
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}
		utils.SendSuccess(w, data)
	}
}
*/

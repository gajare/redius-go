package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gajare/redius-go/models"
	"github.com/gajare/redius-go/service"

	"github.com/gorilla/mux"
)

func CreateTeacher(w http.ResponseWriter, r *http.Request) {
	var t models.Teacher
	_ = json.NewDecoder(r.Body).Decode(&t)
	err := service.CreateTeacher(t)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode("Teacher created")
}

func GetTeacher(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	t, err := service.GetTeacher(id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	json.NewEncoder(w).Encode(t)
}

func UpdateTeacher(w http.ResponseWriter, r *http.Request) {
	var t models.Teacher
	_ = json.NewDecoder(r.Body).Decode(&t)
	err := service.UpdateTeacher(t)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode("Teacher updated")
}

func DeleteTeacher(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := service.DeleteTeacher(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode("Teacher deleted")
}

//multiple teacher inserted

func CreateTeachers(w http.ResponseWriter, r *http.Request) {
	var teachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&teachers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, t := range teachers {
		if err := service.CreateTeacher(t); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode("All teachers inserted")
}

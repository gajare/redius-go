package router

import (
	"net/http"

	"github.com/gajare/redius-go/controller"

	"github.com/gorilla/mux"
)

func Router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/teacher", controller.CreateTeacher).Methods("POST")
	r.HandleFunc("/teacher/{id}", controller.GetTeacher).Methods("GET")
	r.HandleFunc("/teacher", controller.UpdateTeacher).Methods("PUT")
	r.HandleFunc("/teacher/{id}", controller.DeleteTeacher).Methods("DELETE")
	r.HandleFunc("/teachers", controller.CreateTeachers).Methods("POST")

	return r
}

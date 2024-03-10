package router

import (
	"github.com/gorilla/mux"
	Controller "github.com/shubash/pipo/controller"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/", Controller.Home).Methods("GET")
	router.HandleFunc("/api/singup", Controller.Singup).Methods("POST")
	router.HandleFunc("/api/login", Controller.Login).Methods("POST")
	router.HandleFunc("/api/users", Controller.GetAllUser).Methods("GET")
	router.HandleFunc("/api/update/{id}", Controller.UPDATEPASSWORD).Methods("PUT")

	router.HandleFunc("/api/deleteusers", Controller.DeleteAll).Methods("DELETE")

	return router
}

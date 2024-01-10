package router

import (
	"github.com/gorilla/mux"
	"libraryManagementSystem/controllers"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/books/newbook", controllers.CreateBook).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/books/getbook/{id}", controllers.GetBook).Methods("GET", "OPTIONS")
	return router
}

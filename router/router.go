package router

import(
	"libraryManagementSystem/controllers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router:= mux.NewRouter();
	router.HandleFunc("/api/books/newbook", controllers.CreateBook).Methods("POST", "OPTIONS")
	return router
}
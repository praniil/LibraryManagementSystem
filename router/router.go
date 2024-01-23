package router

import (
	"github.com/gorilla/mux"
	"libraryManagementSystem/controllers"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/books/newbook", controllers.CreateBook).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/books/getbook/{id}", controllers.GetBook).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/books/getallbooks", controllers.GetAllBook).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/books/updatebook/{id}", controllers.UpdateBook).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/books/deletebook/{id}", controllers.DeleteBook).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/students/newstudent", controllers.CreateStudent).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/students/getstudent/{id}", controllers.GetStudent).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/students/getallstudents", controllers.GetAllStudents).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/students/updatestudent/{id}", controllers.UpdateStudent).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/students/deletestudent/{id}", controllers.DeleteStudent).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/loans/returnbook", controllers.ReturnBook).Methods("POST", "OPTIONS")
	return router
}

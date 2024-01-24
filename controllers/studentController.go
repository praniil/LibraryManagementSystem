package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"libraryManagementSystem/database"
	"libraryManagementSystem/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func CreateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var student models.StudentInformation
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		log.Fatalf("failed to decode json format to the original format. %v", err)
	}
	bookId := student.BooksId
	studentId := insertStudent(student, bookId)
	var msg string
	if studentId == 0 {
		msg = "Student not created"
	} else {
		msg = "A new Student crested"
	}
	res := Response{
		ID:      studentId,
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func insertStudent(student models.StudentInformation, bookId int) int64 {
	db := database.Database_connection()
	// if exists := db.Migrator().HasTable(&models.BookInfo{}); exists {
	// 	fmt.Println("Table books exists")
	// } else {
	// 	db.AutoMigrate(&models.BookInfo{}, &models.StudentInfo{})
	// }
	db.AutoMigrate(&models.BookInformation{}, &models.StudentInformation{}, &models.LoanInformation{})
	var book models.BookInformation
	tx := db.Begin()

	tx.Model(&models.BookInformation{}).Where("id = ? AND students_id = ?", bookId, 0).Find(&book)

	if book.ID == 0 { //no book found with title given by student and the book is not available
		tx.Rollback()
		log.Fatalf("Book not found for title: %d and students_id: %d", student.BooksId, 0)
	}

	result := tx.Create(&student)

	if result.Error != nil {
		tx.Rollback()
		log.Fatalf("Unable to create a record for students. %v", result.Error)
	}

	book.StudentsId = int(student.ID)
	book.StudentsFullName = student.FullName

	tx.Model(&models.BookInformation{}).Where("id =?", book.ID).Updates(&book)

	var loanDuration time.Duration
	loanDuration = 2 * 7 * 24 * time.Hour
	dueDate := time.Now().Add(loanDuration)
	remainingTime := dueDate.Sub(time.Now())
	loanInfos := models.LoanInformation{
		BookID:           int64(book.ID),
		StudentsID:       int64(student.ID),
		DueDate:          dueDate,
		RemainingTime:    remainingTime,
		Fine:             0,
		Returned:         false,
		StudentsFullName: student.FullName,
	}
	result = tx.Create(&loanInfos)
	if result.Error != nil {
		tx.Rollback()
		log.Fatalf("Unable to create a loan record. %v", result.Error)
	}
	var loans []models.LoanInformation
	resultUpdateLoan := tx.Find(&loans)
	if resultUpdateLoan.Error != nil {
		fmt.Printf("error finding the loans information. %v", resultUpdateLoan.Error)
	}
	for i := range loans {
		loans[i].RemainingTime = loans[i].DueDate.Sub(time.Now())
		if loans[i].RemainingTime < 0 {
			loans[i].RemainingTime = 0
		}
	}
	tx.Save(&loans)
	tx.Commit()
	return int64(student.ID)

}

func GetStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to extract id from the url. %v", err)
	}

	student, err := getStudent(int64(id))
	if err != nil {
		fmt.Println("unable to get the student")
	}
	json.NewEncoder(w).Encode(student)
}

func getStudent(id int64) (models.StudentInformation, error) {
	db := database.Database_connection()
	var student models.StudentInformation
	result := db.First(&student, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Printf("couldnt find the record with %d, %v", id, result.Error)
		return student, result.Error
	}
	if result.Error != nil {
		fmt.Printf("couldnt get the record with id : %d, error: %v", id, result.Error)
		return student, result.Error
	}
	return student, nil
}
func GetAllStudents(w http.ResponseWriter, r *http.Request) {
	//r *http.Request http.Request struct contains information about an incomming HTTP request from a client, using a pointer to this struct allows the function to access and modify the req data
	//its like func insert(node *Node)
	w.Header().Set("Content-type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	students, err := getAllStudents()
	if err != nil {
		log.Fatalf("unable to find students. %v", err)
	}

	json.NewEncoder(w).Encode(students)
}

func getAllStudents() ([]models.StudentInformation, error) {
	db := database.Database_connection()
	var students []models.StudentInformation //{arrays of information of different students}
	result := db.Find(&students)             //retrieves all the information from models.Student table

	if result.Error != nil {
		log.Fatalf("unable to find students. %v", result.Error)
	}
	return students, nil
}

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"]) //{id} extract
	if err != nil {
		log.Fatalf("couldnot extract id from the url, %v", err)
	}
	var student models.StudentInformation
	json.NewDecoder(r.Body).Decode(&student)

	rowsUpdated := updateStudent(student, int64(id))
	msg := fmt.Sprintf("number of rows updated: %d", rowsUpdated)
	res := Response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)

}

func updateStudent(student models.StudentInformation, id int64) int64 {
	db := database.Database_connection()
	result := db.Model(&models.StudentInformation{}).Where("id = ?", id).Updates(student)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Printf("Record not found with id: %d", id)
	}
	if result.Error != nil {
		log.Fatalf("unable to update the record . %v", result.Error)
	}
	rowsUpdated := result.RowsAffected
	return rowsUpdated

}

func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("couldnt extract id from the url. %v", err)
	}

	deletedId := deleteStudent(int64(id))
	res := Response{
		ID:      deletedId,
		Message: "this row is deleted",
	}
	json.NewEncoder(w).Encode(res)
}

func deleteStudent(id int64) int64 {
	db := database.Database_connection()
	result := db.Delete(&models.BookInformation{}, id)
	if result.Error != nil {
		log.Fatalf("couldnt delete the row %v", result.Error)
	}
	rowsDeleted := result.RowsAffected
	return rowsDeleted
}

//create a way to return book
/*
func ReturnBook(loanID uint) {
    var loan models.Loan
    db.First(&loan, loanID)

    if loan.ID == 0 {
        fmt.Println("Loan not found")
        return
    }

    if loan.Returned {
        fmt.Println("Book has already been returned")
        return
    }

    // Calculate the days overdue
    daysOverdue := int(math.Max(0, time.Since(loan.DueDate).Hours()/24))

    // Calculate the fine
    fineAmount := finePerDay * float64(daysOverdue)

    // Update the loan information
    loan.Returned = true
    loan.Fine = fineAmount
    db.Save(&loan)

    fmt.Printf("Book with ID %d returned. Fine amount: $%.2f\n", loan.BookID, fineAmount)
}
*/

//its the idea but we make a new api for this shit

func ReturnBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var student models.StudentInformation
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		fmt.Printf("couldnt convert the json format to the original format, %v", err)
	}

	loanID := returnBook(student, student.BooksId)
	msg := fmt.Sprintf("The book is returned and the loan is cleared of id: %d", loanID)
	res := Response{
		ID:      loanID,
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func returnBook(student models.StudentInformation, bookId int) int64 {
	db := database.Database_connection()
	var loan models.LoanInformation
	result := db.Where("book_id = ? AND students_full_name = ?", bookId, student.FullName).Find(&loan)
	if result.Error != nil {
		fmt.Printf("couldnot find the loan information with bookid:%d error: %v", bookId, result.Error)
	}
	loan.RemainingTime = 0
	loan.Returned = true
	// updatedFields := map[string]interface{}{
	// 	"remaining_time": loan.RemainingTime,
	// 	"returned":       loan.Returned,
	// }
	db.Model(&loan).Where("book_id =?", bookId).Update("remaining_time", loan.RemainingTime)
	db.Model(&loan).Where("book_id =?", bookId).Update("returned", loan.Returned)
	fmt.Println("loan id:", loan.ID)
	fmt.Println("loan full name", loan.StudentsFullName)
	return int64(loan.ID)
}

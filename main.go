package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Person struct {
	gorm.Model

	Name  string
	Email string `gorm:"typevarchar(100);unique_index"`
	Book  []Book
}

type Book struct {
	gorm.Model

	Title      string
	Author     string
	CallNumber int `gorm:"unique_index"`
	PersonID   int
}

var db *gorm.DB
var err error

// var (
// 	person = &Person{Name: "Rohan", Email: "rohan@email.com"}
// 	books  = []Book{{Title: "Book 1", Author: "Author 1", CallNumber: 123, PersonID: 2},
// 		{Title: "Book 2", Author: "Author 2", CallNumber: 1234, PersonID: 3}}
// )

func main() {

	//*loading ENv variable
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	port := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("NAME")
	passwd := os.Getenv("PASSWORD")

	DBURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbname, passwd, port)
	//* Starting/Opening Database connection via gorm
	db, err = gorm.Open(dialect, DBURI)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("DB Connected Successfully!")
	}
	//* closing when done
	defer db.Close()

	//* Make migration to the database if they have not already been created
	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Book{})

	//Mux Router

	router := mux.NewRouter()
	fmt.Println("API running")

	router.HandleFunc("/people", GetPeople).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	var people []Person
	db.Find(&people)
	json.NewEncoder(w).Encode(&people)
}

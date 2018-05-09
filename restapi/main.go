package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var books []Book

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

func createBook(w http.ResponseWriter, r *http.Request) {

}

func updateBook(w http.ResponseWriter, r *http.Request) {

}

func deleteBook(w http.ResponseWriter, r *http.Request) {

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	r := mux.NewRouter()

	db, err := sql.Open("sqlite3", "./foo.db")
	checkErr(err)

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS books (id INTEGER PRIMARY KEY, isbn TEXT, title TEXT)")
	stmt.Exec()

	//delete

	stmt, err = db.Prepare("delete from books")
	checkErr(err)

	res, err := stmt.Exec()
	checkErr(err)

	// insert
	stmt, err = db.Prepare("INSERT INTO books(id, isbn, title) values(?,?,?)")
	checkErr(err)

	res, err = stmt.Exec("1", "xxxxx", "book1")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)

	rows, err := db.Query("SELECT * FROM books")
	checkErr(err)

	columns, err := rows.Columns()
	checkErr(err)

	for _, element := range columns {
		fmt.Println(element)
	}

	var book Book

	for rows.Next() {
		err = rows.Scan(&book.ID, &book.Isbn, &book.Title)
		checkErr(err)
	}

	rows.Close() //good habit to close

	books = append(books, book)

	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}

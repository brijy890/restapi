package main

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/brij/restapi/db"
)

// Book struct (MOdel)

type Book struct {
	Isbn  string `json:"isbn"`
	Title string `json:"title"`
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	db := db.CreatedbConn()
	w.Header().Set("Content-Type", "application/json")
	selDB, err := db.Query("SELECT isbn, title FROM book ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	books := Book{}
	res := []Book{}
	for selDB.Next() {
		var isbn, title string
		err = selDB.Scan(&isbn, &title)
		if err != nil {
			panic(err.Error())
		}
		books.Isbn = isbn
		books.Title = title
		res = append(res, books)
	}
	json.NewEncoder(w).Encode(res)
	defer db.Close()
}

func getBook(w http.ResponseWriter, r *http.Request) {
	db := db.CreatedbConn()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	selDB, err := db.Query("SELECT isbn, title FROM book WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	books := Book{}
	res := []Book{}
	for selDB.Next() {
		var isbn, title string
		err = selDB.Scan(&isbn, &title)
		if err != nil {
			panic(err.Error())
		}
		books.Isbn = isbn
		books.Title = title
		res = append(res, books)
	}
	json.NewEncoder(w).Encode(res)
	defer db.Close()
}

func createBook(w http.ResponseWriter, r *http.Request) {
	db := db.CreatedbConn()
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	insForm, err := db.Prepare("INSERT INTO book(isbn, title) VALUES(?,?)")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(book.Isbn, book.Title)
	json.NewEncoder(w).Encode(book)
	defer db.Close()
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	db := db.CreatedbConn()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	id := params["id"]
	insForm, err := db.Prepare("UPDATE book SET isbn=?, title=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(book.Isbn, book.Title, id)
	defer db.Close()
	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	db := db.CreatedbConn()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	delForm, err := db.Prepare("DELETE FROM book WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(id)
	json.NewEncoder(w).Encode(id)
	defer db.Close()
}

func main() {
	// Init router
	r := mux.NewRouter()

	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Println("Server started on: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

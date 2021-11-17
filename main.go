package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Category struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

var db *sql.DB
var err error

func main() {

	db, err = sql.Open("mysql", "root:my-secret@tcp(127.0.0.1:3306)/sample")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/categories", getCategories).Methods("GET")
	router.HandleFunc("/categories", createCategory).Methods("POST")
	router.HandleFunc("/categories/{id}", getCategory).Methods("GET")
	router.HandleFunc("/categories/{id}", updateCategory).Methods("PUT")
	router.HandleFunc("/categories/{id}", deleteCategory).Methods("DELETE")

	http.ListenAndServe(":8000", router)
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var categories []Category

	result, err := db.Query("SELECT id, title from categories")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var category Category
		err := result.Scan(&category.ID, &category.Title)
		if err != nil {
			panic(err.Error())
		}
		categories = append(categories, category)
	}

	json.NewEncoder(w).Encode(categories)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("INSERT INTO categories(title) VALUES(?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["title"]
	_, err = stmt.Exec(title)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "New category was created")
}
func getCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, title FROM categories WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var category Category
	for result.Next() {
		err := result.Scan(&category.ID, &category.Title)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(category)
}
func updateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE categories SET title = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newTitle := keyVal["title"]
	_, err = stmt.Exec(newTitle, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Category with ID = %s was updated", params["id"])
}
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM categories WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Category with ID = %s was deleted", params["id"])
}

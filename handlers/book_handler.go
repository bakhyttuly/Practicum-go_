package handlers

import (
	"bookstore/config"
	"bookstore/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	rows, err := config.DB.Query("SELECT id, title, author_id, category_id, price FROM books LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var booksList []models.Book
	for rows.Next() {
		var b models.Book
		rows.Scan(&b.ID, &b.Title, &b.AuthorID, &b.CategoryID, &b.Price)
		booksList = append(booksList, b)
	}

	if booksList == nil {
		booksList = []models.Book{}
	}
	json.NewEncoder(w).Encode(booksList)
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := config.DB.QueryRow(
		"INSERT INTO books (title, author_id, category_id, price) VALUES ($1, $2, $3, $4) RETURNING id",
		book.Title, book.AuthorID, book.CategoryID, book.Price,
	).Scan(&book.ID)

	if err != nil {
		http.Error(w, "Could not save book: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func GetBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var b models.Book
	err := config.DB.QueryRow("SELECT id, title, author_id, category_id, price FROM books WHERE id = $1", id).
		Scan(&b.ID, &b.Title, &b.AuthorID, &b.CategoryID, &b.Price)

	if err == sql.ErrNoRows {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(b)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var b models.Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(
		"UPDATE books SET title=$1, author_id=$2, category_id=$3, price=$4 WHERE id=$5",
		b.Title, b.AuthorID, b.CategoryID, b.Price, id,
	)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	b.ID = id
	json.NewEncoder(w).Encode(b)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	_, err := config.DB.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

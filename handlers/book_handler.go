package handlers

import (
	"bookstore/models"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
)

var books = make(map[int]models.Book)
var nextBookID = 1

func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var booksList []models.Book
	for _, book := range books {
		booksList = append(booksList, book)
	}

	sort.Slice(booksList, func(i, j int) bool {
		return booksList[i].ID < booksList[j].ID
	})

	page := 1
	limit := 10

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			http.Error(w, "invalid page", http.StatusBadRequest)
			return
		}
		page = p
	}

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 1 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = l
	}

	start := (page - 1) * limit
	end := start + limit

	if start >= len(booksList) {
		json.NewEncoder(w).Encode([]models.Book{})
		return
	}

	if end > len(booksList) {
		end = len(booksList)
	}

	json.NewEncoder(w).Encode(booksList[start:end])
}

func GetBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	book, exists := books[id]
	if !exists {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(book)
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if book.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	if book.AuthorID <= 0 {
		http.Error(w, "author_id is required", http.StatusBadRequest)
		return
	}
	if book.CategoryID <= 0 {
		http.Error(w, "category_id is required", http.StatusBadRequest)
		return
	}
	if book.Price <= 0 {
		http.Error(w, "price must be greater than 0", http.StatusBadRequest)
		return
	}

	book.ID = nextBookID
	nextBookID++
	books[book.ID] = book

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	_, exists := books[id]
	if !exists {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	var updatedBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updatedBook.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	if updatedBook.AuthorID <= 0 {
		http.Error(w, "author_id is required", http.StatusBadRequest)
		return
	}
	if updatedBook.CategoryID <= 0 {
		http.Error(w, "category_id is required", http.StatusBadRequest)
		return
	}
	if updatedBook.Price <= 0 {
		http.Error(w, "price must be greater than 0", http.StatusBadRequest)
		return
	}

	updatedBook.ID = id
	books[id] = updatedBook

	json.NewEncoder(w).Encode(updatedBook)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	_, exists := books[id]
	if !exists {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	delete(books, id)
	w.WriteHeader(http.StatusNoContent)
}

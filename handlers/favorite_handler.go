package handlers

import (
	"bookstore/config"
	"bookstore/middleware"
	"bookstore/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetFavorites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	page := 1
	limit := 10

	if p := r.URL.Query().Get("page"); p != "" {
		val, err := strconv.Atoi(p)
		if err != nil || val < 1 {
			http.Error(w, "invalid page", http.StatusBadRequest)
			return
		}
		page = val
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		val, err := strconv.Atoi(l)
		if err != nil || val < 1 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = val
	}

	offset := (page - 1) * limit

	rows, err := config.DB.Query(`
        SELECT b.id, b.title, b.author_id, b.category_id, b.price
        FROM favorite_books f
        JOIN books b ON b.id = f.book_id
        WHERE f.user_id = $1
        ORDER BY f.created_at DESC
        LIMIT $2 OFFSET $3
    `, userID, limit, offset)
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.AuthorID, &b.CategoryID, &b.Price); err != nil {
			http.Error(w, "scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, b)
	}

	if result == nil {
		result = []models.Book{}
	}

	json.NewEncoder(w).Encode(result)
}

func AddToFavorites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	bookID, err := strconv.Atoi(params["bookId"])
	if err != nil || bookID < 1 {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	var exists bool
	err = config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM books WHERE id = $1)", bookID).Scan(&exists)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	_, err = config.DB.Exec(`
        INSERT INTO favorite_books (user_id, book_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, book_id) DO NOTHING
    `, userID, bookID)
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "added to favorites"})
}

func RemoveFromFavorites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	bookID, err := strconv.Atoi(params["bookId"])
	if err != nil || bookID < 1 {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(`
        DELETE FROM favorite_books
        WHERE user_id = $1 AND book_id = $2
    `, userID, bookID)
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "favorite not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

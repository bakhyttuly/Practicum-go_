package handlers

import (
	"bookstore/models"
	"encoding/json"
	"net/http"
)

var nextAuthorID = 1
var authors []models.Author

func GetAuthors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authors)
}

func AddAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var author models.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if author.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	author.ID = nextAuthorID
	nextAuthorID++
	authors = append(authors, author)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(author)
}

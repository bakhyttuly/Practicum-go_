package handlers

import (
	"bookstore/models"
	"encoding/json"
	"net/http"
)

var categories []models.Category
var nextCategoryID = 1

func GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func AddCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if category.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	category.ID = nextCategoryID
	nextCategoryID++
	categories = append(categories, category)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

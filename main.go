package main

import (
	"bookstore/config"
	"bookstore/middleware"
	"log"
	"net/http"

	"bookstore/handlers"

	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDB()

	r := mux.NewRouter()

	r.HandleFunc("/auth/token", handlers.GenerateToken).Methods("GET")

	protected := r.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/books/favorites", handlers.GetFavorites).Methods("GET")
	protected.HandleFunc("/books/{bookId}/favorites", handlers.AddToFavorites).Methods("PUT")
	protected.HandleFunc("/books/{bookId}/favorites", handlers.RemoveFromFavorites).Methods("DELETE")

	r.HandleFunc("/books", handlers.GetBooks).Methods("GET")
	r.HandleFunc("/books", handlers.AddBook).Methods("POST")

	r.HandleFunc("/books/{id}", handlers.GetBookByID).Methods("GET")
	r.HandleFunc("/books/{id}", handlers.UpdateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", handlers.DeleteBook).Methods("DELETE")

	r.HandleFunc("/authors", handlers.GetAuthors).Methods("GET")
	r.HandleFunc("/authors", handlers.AddAuthor).Methods("POST")
	r.HandleFunc("/categories", handlers.GetCategories).Methods("GET")
	r.HandleFunc("/categories", handlers.AddCategory).Methods("POST")

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

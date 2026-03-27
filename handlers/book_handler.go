package handlers

import (
	"bookstore/models"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

var books = make(map[int]models.Book)
var nextBookID = 1

func GetBooks(c *gin.Context) {
	var booksList []models.Book
	for _, book := range books {
		booksList = append(booksList, book)
	}

	sort.Slice(booksList, func(i, j int) bool {
		return booksList[i].ID < booksList[j].ID
	})

	page := 1
	limit := 10

	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
			return
		}
		page = p
	}

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}
		limit = l
	}

	start := (page - 1) * limit
	end := start + limit

	if start >= len(booksList) {
		c.JSON(http.StatusOK, []models.Book{})
		return
	}

	if end > len(booksList) {
		end = len(booksList)
	}

	c.JSON(http.StatusOK, booksList[start:end])
}

func GetBookByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	book, exists := books[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func AddBook(c *gin.Context) {
	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if book.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}
	if book.AuthorID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "author_id is required"})
		return
	}
	if book.CategoryID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category_id is required"})
		return
	}
	if book.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price must be greater than 0"})
		return
	}

	book.ID = nextBookID
	nextBookID++
	books[book.ID] = book

	c.JSON(http.StatusCreated, book)
}

func UpdateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	_, exists := books[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	var updatedBook models.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updatedBook.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}
	if updatedBook.AuthorID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "author_id is required"})
		return
	}
	if updatedBook.CategoryID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category_id is required"})
		return
	}
	if updatedBook.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price must be greater than 0"})
		return
	}

	updatedBook.ID = id
	books[id] = updatedBook

	c.JSON(http.StatusOK, updatedBook)
}

func DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	_, exists := books[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	delete(books, id)
	c.Status(http.StatusNoContent)
}

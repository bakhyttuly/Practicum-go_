package handlers

import (
	"bookstore/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var nextAuthorID = 1
var authors []models.Author

func GetAuthors(c *gin.Context) {
	c.JSON(http.StatusOK, authors)
}

func AddAuthor(c *gin.Context) {
	var author models.Author

	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if author.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	author.ID = nextAuthorID
	nextAuthorID++
	authors = append(authors, author)

	c.JSON(http.StatusCreated, author)
}

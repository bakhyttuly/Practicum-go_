package handlers

import (
	"bookstore/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var categories []models.Category
var nextCategoryID = 1

func GetCategories(c *gin.Context) {
	c.JSON(http.StatusOK, categories)
}

func AddCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if category.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	category.ID = nextCategoryID
	nextCategoryID++
	categories = append(categories, category)

	c.JSON(http.StatusCreated, category)
}

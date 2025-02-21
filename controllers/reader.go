package controllers

import (
	"lms/backend/initializers"
	"lms/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchBooks(c *gin.Context) {
	// var book models.BookInventory

	var books []models.BookInventory //will store all the books in the array
	query := c.Query("q")
	if err := initializers.DB.Where("title LIKE ? OR authors LIKE ? OR publisher LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&books).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"Error":   err.Error(),
			"Message": "No book found",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"Books": books,
	})

}

// // SearchBooks enables searching by title, authors, or publisher.
// func SearchBooks(c *gin.Context) {
//     query := c.Query("q")
//     var books []models.Book
//     config.DB.Where("title ILIKE ? OR authors ILIKE ? OR publisher ILIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&books)
//     c.JSON(http.StatusOK, books)

// db.Where("name LIKE ?", "%jin%").Find(&users)
// SELECT * FROM users WHERE name LIKE '%jin%';
// }

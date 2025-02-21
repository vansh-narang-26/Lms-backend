package controllers

import (
	"lms/backend/initializers"
	"lms/backend/models"
	"net/http"
	"time"

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
		return
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

func RaiseIssueRequest(c *gin.Context) {
	userId, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	isbn := c.Param("id")
	email, _ := c.Get("email")
	// fmt.Println("Book id ", isbn)

	//finding the id for the person who has logged in
	var user models.User
	if err := initializers.DB.Where("id = ? AND role = ?", userId, "reader").First(&user).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var book models.BookInventory
	//Find the book with this above id
	if err := initializers.DB.Where("isbn=?", isbn).Find(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error":   err.Error(),
			"Message": "Book not with this isbn",
		})
	}

	//If book not available the request must get rejected
	if book.AvailableCopies <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "Book is not available you request has been denied",
		})
		return
	}
	// var request models.RequestEvent
	// t:=time.Now
	crequest := models.RequestEvent{
		BookID:      isbn,
		ReaderID:    user.ID,
		RequestDate: time.Now(),
		RequestType: "Requested",
	}
	initializers.DB.Create(&crequest)

	c.JSON(http.StatusCreated, gin.H{
		"Message": "Request created successfully",
		"Request": crequest,
		"Email":   email,
	})
}

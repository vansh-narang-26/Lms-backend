package controllers

import (
	"lms/backend/initializers"
	"lms/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddBooks struct {
	ISBN      string
	Title     string
	Author    string
	Publisher string
	Version   int
}

// func AddBook(c *gin.Context) {
// 	var addBook AddBooks
// 	var exisitingUser models.BookInventory

// 	if err := c.ShouldBindJSON(&addBook); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"Error": err.Error(),
// 		})
// 		return
// 	}
// 	//var existingBook models.BookInventory

// 	// if already existed then increase the count
// 	// res := initializers.DB.Where("title=?", existingBook.Title).Find(&existingBook)

// 	// if errors.Is(res.Error, gorm.ErrRecordNotFound) {
// 	//record not found
// 	res := initializers.DB.Where("title=?", addBook.Title).First(&exisitingUser)

// 	// fmt.Println(res)
// 	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
// 		// fmt.Println("record not found")
// 		newBook := models.BookInventory{
// 			ISBN:            addBook.ISBN,
// 			LibID:           1,
// 			Title:           addBook.Title,
// 			Authors:         addBook.Author,
// 			Publisher:       addBook.Publisher,
// 			Version:         addBook.Version,
// 			TotalCopies:     24,
// 			AvailableCopies: 12,
// 		}
// 		initializers.DB.Create(&newBook)
// 		c.JSON(http.StatusOK, gin.H{"data": newBook})
// 	} else {
// 		fmt.Println("record found")
// 		// copies := exisitingUser.TotalCopies
// 		initializers.DB.Model(&models.BookInventory{}).Where("title", addBook.Title).Update("TotalCopies", exisitingUser.TotalCopies+1)
// 		// initializers.DB.Model(&models.BookInventory{})
// 	}
// 	// // fmt.Println("Record not found")

// }

// func RemoveBook(c *gin.Context) {

// }
func AddBook(c *gin.Context) {
	//taking email to return to frontend to see which admin made is creating the book
	email, _ := c.Get("email")
	adminID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	//Btw already being checked in the frontend but checking admin access again
	var adminUser models.User
	if err := initializers.DB.Where("id = ? AND role = ?", adminID, "admin").First(&adminUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	//Checking the json format
	var bookInput models.BookInventory
	if err := c.ShouldBindJSON(&bookInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Checking book exists with the same ISBN or not if exists increase the count by 1
	var existingBook models.BookInventory
	if err := initializers.DB.Where("isbn = ?", bookInput.ISBN).First(&existingBook).Error; err == nil {
		// If book exists, increase the total copies count
		existingBook.TotalCopies += 1
		existingBook.AvailableCopies += 1
		if err := initializers.DB.Save(&existingBook).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book copies"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Book copies updated successfully", "book": existingBook, "email": email})
		return
	}

	//Book not found (else case)
	newBook := models.BookInventory{
		ISBN:            bookInput.ISBN,
		Title:           bookInput.Title,
		Authors:         bookInput.Authors,
		Publisher:       bookInput.Publisher,
		Version:         bookInput.Version,
		TotalCopies:     1,
		AvailableCopies: 1,
		LibID:           adminUser.LibID,
	}

	//Book created and now adding to the DB
	if err := initializers.DB.Create(&newBook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book"})
		return
	}

	//Added the book
	c.JSON(http.StatusCreated, gin.H{"message": "Book added successfully", "book": newBook, "EmailId": email})
}
func RemoveBook(c *gin.Context) {
	isbn := c.Param("id")

	var book models.BookInventory
	if err := initializers.DB.Where("isbn = ?", isbn).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found."})
		return
	}
	if book.AvailableCopies <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No available copies to remove."})
		return
	}
	book.TotalCopies -= 1
	book.AvailableCopies -= 1
	initializers.DB.Save(&book)
	c.JSON(http.StatusOK,
		gin.H{"book": book, "Message": "Book deleted Successfully"})
}

type UpdatorBook struct {
	ISBN      string
	Title     string
	Authors   string
	Publisher string
	Version   int
}

func UpdateBook(c *gin.Context) {
	isbn := c.Param("id")
	// fmt.Println(isbn)
	var upBook UpdatorBook

	if err := c.ShouldBindJSON(&upBook); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	var book models.BookInventory
	if err := initializers.DB.Where("isbn = ?", isbn).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found."})
		return
	}

	book.Title = upBook.Title
	book.Authors = upBook.Authors
	book.Publisher = upBook.Publisher
	book.Version = upBook.Version
	initializers.DB.Save(&book)
	c.JSON(http.StatusOK, book)
}


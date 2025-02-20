package controllers

import (
	"errors"
	"fmt"
	"lms/backend/initializers"
	"lms/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AddBooks struct {
	ISBN      string
	Title     string
	Author    string
	Publisher string
	Version   int
}

func AddBook(c *gin.Context) {
	var addBook AddBooks
	var exisitingUser models.BookInventory

	if err := c.ShouldBindJSON(&addBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}
	//var existingBook models.BookInventory

	// if already existed then increase the count
	// res := initializers.DB.Where("title=?", existingBook.Title).Find(&existingBook)

	// if errors.Is(res.Error, gorm.ErrRecordNotFound) {
	//record not found
	res := initializers.DB.Where("title=?", addBook.Title).First(&exisitingUser)

	// fmt.Println(res)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// fmt.Println("record not found")
		newBook := models.BookInventory{
			ISBN:            addBook.ISBN,
			LibID:           1,
			Title:           addBook.Title,
			Authors:         addBook.Author,
			Publisher:       addBook.Publisher,
			Version:         addBook.Version,
			TotalCopies:     24,
			AvailableCopies: 12,
		}
		initializers.DB.Create(&newBook)
		c.JSON(http.StatusOK, gin.H{"data": newBook})
	} else {
		fmt.Println("record found")
		// copies := exisitingUser.TotalCopies
		initializers.DB.Model(&models.BookInventory{}).Where("title", addBook.Title).Update("TotalCopies", exisitingUser.TotalCopies+1)
		// initializers.DB.Model(&models.BookInventory{})
	}
	// // fmt.Println("Record not found")

}

// func RemoveBook(c *gin.Context) {

// }

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
	c.JSON(http.StatusOK, book)
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

// // SearchBooks enables searching by title, authors, or publisher.
// func SearchBooks(c *gin.Context) {
//     query := c.Query("q")
//     var books []models.Book
//     config.DB.Where("title ILIKE ? OR authors ILIKE ? OR publisher ILIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&books)
//     c.JSON(http.StatusOK, books)
// }

package controllers

import (
	"fmt"
	"lms/backend/initializers"
	"lms/backend/models"
	"net/http"
	"time"

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
	if err := initializers.DB.Where("isbn = ? AND lib_id=?", bookInput.ISBN, adminUser.LibID).First(&existingBook).Error; err == nil {
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

	//check for same isbn doesnt exist
	fmt.Println(bookInput.ISBN)
	if err := initializers.DB.Where("isbn=?", bookInput.ISBN).First(&existingBook).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Same ISBN already exists, cannot add book."})
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
	fmt.Println(adminUser.LibID)

	//Book created and now adding to the DB
	if err := initializers.DB.Create(&newBook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//Added the book
	c.JSON(http.StatusCreated, gin.H{"message": "Book added successfully", "book": newBook, "email": email})
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
	// if book.TotalCopies <= 1 {
	// 	initializers.DB.Delete(&book)
	// 	c.JSON(http.StatusOK, gin.H{"message": "Book removed completely."})
	// 	return
	// }

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

func ListRequests(c *gin.Context) {
	// Getting the admin ID from the context
	adminID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch the admin details
	var adminUser models.User
	if err := initializers.DB.Where("id = ? AND role = ?", adminID, "admin").First(&adminUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// var onerequest models.RequestEvent
	var requests []models.RequestEvent

	if err := initializers.DB.
		Joins("JOIN book_inventories ON request_events.book_id = book_inventories.isbn").
		Where("book_inventories.lib_id = ? AND request_events.request_type = ?", adminUser.LibID, "Requested").
		Find(&requests).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Could not find the requests",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": requests,
	})
}

type UpdateRequest struct {
	RequestType  string     `json:"request_type"`
	ApprovalDate *time.Time `json:"approval_date,omitempty"`
	ApproverID   *uint      `json:"approver_id,omitempty"`
}

func ApproveRequest(c *gin.Context) {
	adminID, _ := c.Get("id")
	//	fmt.Println(adminID)

	var adminUser models.User
	if err := initializers.DB.Where("id = ? AND role = ?", adminID, "admin").First(&adminUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	//extract the request id and check its availablity or not
	requestId := c.Param("id")
	//fmt.Println(requestId)
	var request models.RequestEvent

	if err := initializers.DB.Where("req_id=?", requestId).Find(&request).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"Message": "Couldnot find the request with this id",
		})
		return
	}
	//time to extract bookId from request id
	bookId := request.BookID
	//fmt.Println(bookId)

	var bookexists models.RequestEvent
	if err := initializers.DB.Where("book_id=?", bookId).Find(&bookexists).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"Message": "Couldnt find the book id with this isbn",
		})
		return
	}

	// readerID := request.ReaderID
	// fmt.Println(readerID)

	var handlereq UpdateRequest

	if err := c.ShouldBindJSON(&handlereq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	//time to update the bookrequest
	bookexists.RequestType = "Issued"
	now := time.Now()
	bookexists.ApprovalDate = &now
	bookexists.ApproverID = &adminUser.ID

	initializers.DB.Save(&bookexists)
	// c.JSON(http.StatusOK, bookexists)
	c.JSON(http.StatusAccepted, gin.H{
		"message":      "updation successfully done",
		"updated book": bookexists,
	})

	//now setup the issue registry accordingly
	// var IssueReg models.IssueRegistry

	issueReg := models.IssueRegistry{
		ISBN:               bookId,
		ReaderID:           request.ReaderID,
		IssueApproverID:    adminUser.ID,
		IssueStatus:        "Issued",
		ExpectedReturnDate: time.Now().AddDate(0, 0, 14), // 2 weeks later

	}
	initializers.DB.Create(&issueReg)

	c.JSON(http.StatusCreated, gin.H{
		"Message":   "Issue registry created successfully",
		"Issue reg": issueReg,
	})

}
func RejectRequest(c *gin.Context) {
	adminID, _ := c.Get("id")
	//	fmt.Println(adminID)

	var adminUser models.User
	if err := initializers.DB.Where("id = ? AND role = ?", adminID, "admin").First(&adminUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	//extract the request id and check its availablity or not
	requestId := c.Param("id")
	//fmt.Println(requestId)
	var request models.RequestEvent

	if err := initializers.DB.Where("req_id=?", requestId).Find(&request).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"Message": "Couldnot find the request with this id",
		})
		return
	}
	//time to extract bookId from request id
	bookId := request.BookID
	//fmt.Println(bookId)

	var bookexists models.RequestEvent
	if err := initializers.DB.Where("book_id=?", bookId).Find(&bookexists).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"Message": "Couldnt find the book id with this isbn",
		})
		return
	}

	// readerID := request.ReaderID
	// fmt.Println(readerID)

	var handlereq UpdateRequest

	if err := c.ShouldBindJSON(&handlereq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	//time to update the bookrequest
	fmt.Println(handlereq.RequestType)
	bookexists.RequestType = "Rejected"
	now := time.Now()
	bookexists.ApprovalDate = &now
	bookexists.ApproverID = &adminUser.ID

	initializers.DB.Save(&bookexists)
	// c.JSON(http.StatusOK, bookexists)
	c.JSON(http.StatusAccepted, gin.H{
		"message":      "updation successfully done",
		"updated book": bookexists,
	})

	// //now setup the issue registry accordingly
	// // var IssueReg models.IssueRegistry

	// issueReg := models.IssueRegistry{
	// 	ISBN:               bookId,
	// 	ReaderID:           request.ReaderID,
	// 	IssueApproverID:    adminUser.ID,
	// 	IssueStatus:        "Rejected",
	// 	ExpectedReturnDate: time.Now().AddDate(0, 0, 14), // 2 weeks later

	// }
	// initializers.DB.Create(&issueReg)

	// c.JSON(http.StatusCreated, gin.H{
	// 	"Message":   "Issue registry created successfully",
	// 	"Issue reg": issueReg,
	// })

}

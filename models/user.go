// Library: (ID, Name)
// Users: (ID, Name, Email, ContactNumber, Role, LibID)
// BookInventory : (ISBN, LibID, Title, Authors, Publisher, Version, TotalCopies, AvailableCopies)
// RequestEvents: (ReqID, BookID, ReaderID, RequestDate, ApprovalDate, ApproverID, RequestType)
// IssueRegistery: (IssueID, ISBN, ReaderID, IssueApproverID, IssueStatus, IssueDate, ExpectedReturnDate, ReturnDate, ReturnApproverID)

package models

import "time"

type User struct {
	ID            uint     `json:"id" gorm:"primary_key"`
	Name          string   `json:"name"`
	Email         string   `json:"email"`
	ContactNumber string   `json:"contact_no"`
	Role          string   `json:"role" binding:"oneof=admin reader owner"`
	Library       *Library `gorm:"foreignKey:LibID"`
	LibID         uint     `json:"lib_id"`
}

type LoginUser struct {
	Email string `json:"email"`
	// Role  string `json:"role" binding:"oneof=admin reader"`
}

type Library struct {
	ID   uint `json:"id" gorm:"primary_key"`
	Name string
}

// type Library struct {
//     ID   uint   `gorm:"primaryKey" json:"id"`
//     Name string `gorm:"unique;not null" json:"name"`
// }

//	type User struct {
//	    ID            uint   `gorm:"primaryKey" json:"id"`
//	    Name          string `gorm:"not null" json:"name"`
//	    Email         string `gorm:"unique;not null" json:"email"`
//	    ContactNumber string `json:"contact_number"`
//	    Role          string `gorm:"not null" json:"role"` // Owner, Admin, or Reader
//	    LibID         uint   `json:"lib_id"`
//	    Library       Library `gorm:"foreignKey:LibID" json:"-"`
//	}
//
//	type BookInventory struct {
//		ISBN            string
//		Library         *Library `gorm:"foreignKey:LibID"`
//		LibID           uint     `json:"lib_id"`
//		Title           string
//		Authors         string
//		Publisher       string
//		Version         string
//		TotalCopies     uint
//		AvailableCopies uint
//	}
type BookInventory struct {
	ISBN            string  `gorm:"primaryKey" json:"isbn"`
	LibID           uint    `json:"lib_id"`
	Title           string  `gorm:"not null" json:"title"`
	Authors         string  `json:"authors"`
	Publisher       string  `json:"publisher"`
	Version         int     `json:"version"`
	TotalCopies     int     `gorm:"not null" json:"total_copies"`
	AvailableCopies int     `gorm:"not null" json:"available_copies"`
	Library         Library `gorm:"foreignKey:LibID" json:"-"`
}

type RequestEvent struct {
	ReqID        uint       `gorm:"primaryKey" json:"req_id"`
	BookID       string     `json:"book_id"`
	ReaderID     uint       `json:"reader_id"`
	RequestDate  time.Time  `json:"request_date"`
	ApprovalDate *time.Time `json:"approval_date,omitempty"`
	ApproverID   *uint      `json:"approver_id,omitempty"`
	RequestType  string     `json:"request_type"` // Issued or Returned
	Status       string     `gorm:"default:'pending'" json:"status"`
}

type IssueRegistry struct {
	IssueID            uint       `gorm:"primaryKey" json:"issue_id"`
	ISBN               string     `json:"isbn"`
	ReaderID           uint       `json:"reader_id"`
	IssueApproverID    uint       `json:"issue_approver_id"`
	IssueStatus        string     `json:"issue_status"` // Issued or Returned
	IssueDate          time.Time  `json:"issue_date"`
	ExpectedReturnDate time.Time  `json:"expected_return_date"`
	ReturnDate         *time.Time `json:"return_date,omitempty"`
	ReturnApproverID   *uint      `json:"return_approver_id,omitempty"`
}

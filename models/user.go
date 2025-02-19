// Library: (ID, Name)
// Users: (ID, Name, Email, ContactNumber, Role, LibID)
// BookInventory : (ISBN, LibID, Title, Authors, Publisher, Version, TotalCopies, AvailableCopies)
// RequestEvents: (ReqID, BookID, ReaderID, RequestDate, ApprovalDate, ApproverID, RequestType)
// IssueRegistery: (IssueID, ISBN, ReaderID, IssueApproverID, IssueStatus, IssueDate, ExpectedReturnDate, ReturnDate, ReturnApproverID)

package models

type User struct {
	ID            int      `json:"id" gorm:"primary_key"`
	Name          string   `json:"name"`
	Email         string   `json:"email"`
	ContactNumber string   `json:"contact_no"`
	Role          string   `json:"role" binding:"oneof=admin reader owner"`
	Library       *Library `gorm:"foreignKey:LibID"`
	LibID         int      `json:"lib_id"`
}

type LoginUser struct {
	Email string `json:"email"`
	// Role  string `json:"role" binding:"oneof=admin reader"`
}

type Library struct {
	ID   int `gorm:"type:uuid;primaryKey"`
	Name string
}

type BookInventory struct {
	ISBN            string
	Library         *Library `gorm:"foreignKey:LibID"`
	LibID           int      `json:"lib_id"`
	Title           string
	Authors         string
	Publisher       string
	Version         string
	TotalCopies     uint
	AvailableCopies uint
}

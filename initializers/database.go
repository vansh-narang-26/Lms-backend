package initializers

import (
	"lms/backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("library.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	// Auto Migrate the schemas
	database.AutoMigrate(&models.Library{}, &models.User{}) // Will become libraries,users in the database
	database.AutoMigrate(&models.BookInventory{})
	database.AutoMigrate(&models.RequestEvent{})
	database.AutoMigrate(&models.IssueRegistry{})

	DB = database
}

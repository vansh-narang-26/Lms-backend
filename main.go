package main

import (
	"lms/backend/controllers"
	"lms/backend/initializers"
	"lms/backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	initializers.ConnectDatabase()

	router := gin.Default()

	publicRoutes := router.Group("/api")
	{

		publicRoutes.POST("/users/register", controllers.CreateUser)
		publicRoutes.POST("/users/login", controllers.LoginUser)
	}
	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middleware.UserRetriveCookie)
	{
		library := protectedRoutes.Group("/library")
		library.Use(middleware.OwnerOnly)
		{
			library.POST("/create", controllers.CreateLibrary)
			library.POST("/create-admin", controllers.CreateAdmin)
		}

		admin := protectedRoutes.Group("/admin")
		admin.Use(middleware.AdminOnly)
		{
			admin.POST("/add-book", controllers.AddBook)
			admin.DELETE("/:id", controllers.RemoveBook) // Only need to
			admin.PUT("/:id", controllers.UpdateBook)    // Only need to update with the id
		}
		reader := protectedRoutes.Group("/reader")
		reader.Use(middleware.ReaderOnly)
		{
			reader.GET("/search-books", controllers.SearchBooks)
		}
	}

	// protectedRoutes := router.Group("/api")
	// protectedRoutes.Use(middleware.AuthMiddleware())
	// {
	// 	// Library management (Owner only)
	// 	library := protected.Group("/library")
	// 	library.Use(middleware.OwnerOnly())
	// 	{
	// 		library.PUT("", controllers.UpdateLibrary)
	// 		library.DELETE("", controllers.DeleteLibrary)
	// 	}

	// 	// User management (Owner/Admin)
	// 	users := protected.Group("/users")
	// 	users.Use(middleware.AdminOnly())
	// 	{
	// 		users.POST("", controllers.CreateUser)
	// 		users.GET("", controllers.GetUsers)
	// 	}

	// 	// Book management (Owner/Admin)
	// 	books := protected.Group("/books")
	// 	{
	// 		// Public book search
	// 		books.GET("", controllers.SearchBooks)

	// 		// Admin only operations
	// 		adminBooks := books.Group("")
	// 		adminBooks.Use(middleware.AdminOnly())
	// 		{
	// 			adminBooks.POST("", controllers.AddBook)
	// 			adminBooks.DELETE("/:isbn", controllers.RemoveBook)
	// 		}
	// 	}

	// 	// Issue management
	// 	issues := protected.Group("/issues")
	// 	{
	// 		// Reader operations
	// 		issues.POST("/request", controllers.RaiseIssueRequest)

	// 		// Admin operations
	// 		adminIssues := issues.Group("")
	// 		adminIssues.Use(middleware.AdminOnly())
	// 		{
	// 			adminIssues.POST("/approve", controllers.ApproveIssueRequest)
	// 		}
	// 	}
	// }

	router.Run(":8000")
}

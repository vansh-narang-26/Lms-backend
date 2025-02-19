package main

import (
	"lms/backend/controllers"
	"lms/backend/initializers"

	"github.com/gin-gonic/gin"
)

func main() {
	initializers.ConnectDatabase()

	router := gin.Default()

	router.POST("/auth/signup", controllers.CreateUser)
	router.POST("/auth/login", controllers.LoginUser)
	router.POST("/auth/getusers", controllers.GetUsers)

	//need to make the middleware for the same
	router.POST("/auth/owner/create-lib", controllers.CreateLibrary)
	// router.POST("/auth/login", controllers.Login)
	// router.GET("/user/profile", middlewares.CheckAuth, controllers.GetUserProfile)
	router.Run(":8000")
}

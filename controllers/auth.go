package controllers

import (
	"fmt"
	"lms/backend/initializers"
	"lms/backend/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey = []byte("SECRET")

func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}
	fmt.Println(user)

	var exisitingUser models.User

	initializers.DB.Where("email=?", user.Email).Find(&exisitingUser)

	if exisitingUser.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "already exists with the same email",
		})
		return
	}
	cuser := models.User{
		Name:          user.Name,
		Email:         user.Email,
		ContactNumber: user.ContactNumber,
		Role:          user.Role,
		LibID:         user.LibID,
	}

	initializers.DB.Create(&cuser)

	c.JSON(http.StatusOK, gin.H{"data": cuser})
}

func LoginUser(c *gin.Context) {
	var luser models.LoginUser

	if err := c.ShouldBindJSON(&luser); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"Error": err.Error(),
		})
		return
	}
	var userFound models.User

	//Finding user with email only but when found return the role also in the response
	initializers.DB.Where("email=?", luser.Email).Find(&userFound)

	if userFound.ID == 0 {
		c.JSON(http.StatusBadGateway, gin.H{
			"Error": "No user found",
		})
		return
	}
	claims := jwt.MapClaims{}
	claims["id"] = userFound.ID
	claims["role"] = userFound.Role
	claims["email"] = userFound.Email

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newToken, _ := token.SignedString(secretKey)

	c.SetCookie("Authorise", newToken, 3600, "", "", false, true)

	c.JSON(200, gin.H{
		"message": "Logged in successfully",
		"token":   newToken,
		"role":    userFound.Role,
	})
}

// func GetUsers(context *gin.Context) {
// 	var user []models.User
// 	err := initializers.DB.Find(&user)
// 	if err != nil {
// 		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
// 		return
// 	}
// 	context.JSON(http.StatusOK, user)
// }

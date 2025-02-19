package controllers

import (
	"lms/backend/initializers"
	"lms/backend/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)
func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

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
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": userFound.Role,
		"id":   userFound.ID,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	// fmt.Println(generateToken)

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
	}

	c.JSON(200, gin.H{
		"message": "Logged in successfully",
		"token":   token,
	})
}

func GetUsers(context *gin.Context) {
	var user []models.User
	err := initializers.DB.Find(&user)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	context.JSON(http.StatusOK, user)
}

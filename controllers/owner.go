package controllers

import (
	"lms/backend/initializers"
	"lms/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Library struct {
	ID   int
	Name string
}

func CreateLibrary(c *gin.Context) {
	//ID
	//name

	//clib is creation of library
	var clib models.Library

	if err := c.ShouldBindJSON(&clib); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}
	var existingLibrary models.Library

	initializers.DB.Where("name=?", clib.Name).Find(&existingLibrary)

	if existingLibrary.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Library with this name already exisits",
		})
		return
	}

	nlibrary := models.Library{
		Name: clib.Name,
	}

	initializers.DB.Create(&nlibrary)

	c.JSON(http.StatusOK, gin.H{
		"data": nlibrary,
	})

}

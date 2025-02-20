package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func UserRetriveCookie(c *gin.Context) {

	valid := ValidateCookie(c)
	if valid == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
		c.Abort()
	} else {
		userId, role, err := RetriveJwtToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
		} else {
			c.Set("id", userId)
		}
		if role != "owner" {
			//fmt.Println(role)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Only owner can create the library"})
			c.Abort()
		}
	}
	c.Next()
}
func RetriveJwtToken(c *gin.Context) (int, string, error) {
	cookie, err := c.Cookie("Authorise")
	if err != nil {
		return 0, "", errors.New("cookie not found")
	}

	fmt.Println("Cookie:", cookie)
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte("SECRET"), nil
	})

	if err != nil {
		fmt.Println("Token parse error:", err)
		return 0, "", err
	}

	fmt.Println(token)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId, okId := claims["id"].(float64)
		role, okRole := claims["role"].(string)
		// fmt.Println(userId, role)
		if !okId || !okRole {
			return 0, "", fmt.Errorf("invalid claims")
		}
		return int(userId), role, nil
	} else {
		return 0, "", fmt.Errorf("invalid token")
	}
}
func ValidateCookie(c *gin.Context) bool {
	cookie, _ := c.Cookie("Authorise")
	// fmt.Println(cookie)
	if cookie == "" {
		fmt.Println("cookie not found")
		return false
	} else {
		return true
	}

}

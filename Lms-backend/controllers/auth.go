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

	if err := c.ShouldBindJSON(&user); err != nil {
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
			"Error": "Same email exists",
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

	if err := c.ShouldBindJSON(&luser); err != nil || luser.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var userFound models.User
	initializers.DB.Where("email = ?", luser.Email).Find(&userFound)

	if userFound.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or user does not exist"})
		return
	}

	claims := jwt.MapClaims{
		"id":    userFound.ID,
		"role":  userFound.Role,
		"email": userFound.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newToken, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
		"token":   newToken,
		"role":    userFound.Role,
	})
}

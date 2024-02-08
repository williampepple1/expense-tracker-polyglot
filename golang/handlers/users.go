package handlers

import (
	"errors"
	"expense-tracker/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtKey = os.Getenv("JWT_SECRET_KEY")

func RegisterUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Keep a lowercase version of the username for checking duplicates and saving
		lowerUsername := strings.ToLower(user.Username)

		// Check if a lowercase username already exists
		var existingUser models.User
		result := db.Where("lower(username) = ?", lowerUsername).First(&existingUser)
		if result.Error == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
			return
		}

		// Check if the error is not a 'record not found' error
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing username"})
			return
		}

		// Continue with password encryption and user creation
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}
		user.Password = string(hashedPassword)

		// Use the lowercase username when saving to the database
		user.Username = lowerUsername

		// Create User
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "id": user.ID})
	}
}

func LoginUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var inputUser models.User

		if err := c.ShouldBindJSON(&inputUser); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := db.Where("username = ?", inputUser.Username).First(&user).Error; err != nil {
			c.JSON(400, gin.H{"error": "Username or password incorrect"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputUser.Password)); err != nil {
			c.JSON(400, gin.H{"error": "Username or password incorrect"})
			return
		}

		expirationTime := time.Now().Add(15 * time.Minute)
		claims := &jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   user.Username,
			Id:        user.ID.String(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(jwtKey))

		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(200, gin.H{"token": tokenString})
	}
}

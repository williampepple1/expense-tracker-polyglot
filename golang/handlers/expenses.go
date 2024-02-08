package handlers

import (
	"expense-tracker/models"

	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Helper function to retrieve and validate user ID and expense
func getUserIDAndExpense(c *gin.Context, db *gorm.DB, expenseId string) (uuid.UUID, models.Expense, bool) {
	var expense models.Expense

	// Retrieve the user ID from the context
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(401, gin.H{"error": "User ID not found"})
		return uuid.Nil, expense, false
	}

	// Assert that userID is of type string
	userStrId, ok := userID.(string)
	if !ok {
		c.JSON(400, gin.H{"error": "User ID is not of type string"})
		return uuid.Nil, expense, false
	}

	// Parse the user ID from string to uuid.UUID
	userUUID, err := uuid.Parse(userStrId)
	if err != nil {
		c.JSON(500, gin.H{"error": "User ID is not a valid UUID"})
		return uuid.Nil, expense, false
	}

	// Find the expense by ID and retrieve it along with its UserID
	if err := db.Where("id = ?", expenseId).First(&expense).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Expense not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to retrieve expense"})
		}
		return uuid.Nil, expense, false
	}

	// Check if the user is authorized to access the expense
	if userUUID != expense.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to access this expense"})
		return uuid.Nil, expense, false
	}

	return userUUID, expense, true
}

func ListExpenses(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var expenses []models.Expense
		if err := db.Find(&expenses).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve expenses"})
			return
		}
		c.JSON(http.StatusOK, expenses)
	}
}

func ListUserExpenses(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the user ID from the URL parameter.
		userIDStr := c.Param("userId")

		// Validate that the userID is a valid UUID
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is not a valid UUID"})
			return
		}

		var expenses []models.Expense
		// Find expenses where the UserID matches the provided UUID.
		if err := db.Where("user_id = ?", userID).Find(&expenses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expenses for the user"})
			return
		}

		c.JSON(http.StatusOK, expenses)
	}
}
func CreateExpense(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var expense models.Expense
		if err := c.ShouldBindJSON(&expense); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Retrieve the user ID from the context
		userID, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "User ID not found"})
			return
		}

		// Assert that userID is of type string
		userStrId, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": "User ID is not of type string"})
			return
		}

		// Parse the user ID from string to uuid.UUID
		id, err := uuid.Parse(userStrId)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": "User ID is not a valid UUID"})
			return
		}

		expense.UserID = id // Assign the user's UUID to the expense's UserID field

		// Create the expense in the DB
		if err := db.Create(&expense).Error; err != nil {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "Failed to create expense"})
			return
		}
		// Now load the user associated with the expense and update the expense.User
		var user models.User
		if err := db.Where("id = ?", expense.UserID).First(&user).Error; err != nil {
			// Handle the error. Perhaps the user does not exist in the DB
			c.JSON(500, gin.H{"error": "Failed to load user data"})
			return
		}

		expense.User = user // Assuming your Expense struct has a User field to hold this data

		// Now expense should have the User data included

		c.JSON(http.StatusCreated, expense)
	}
}

func GetExpense(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var expense models.Expense
		if err := db.Where("id = ?", id).First(&expense).Error; err != nil {
			c.JSON(404, gin.H{"error": "Expense not found"})
			return
		}

		c.JSON(http.StatusOK, expense)
	}
}

func UpdateExpense(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Use the helper function to get the user ID and expense
		userUUID, expense, ok := getUserIDAndExpense(c, db, id)
		if !ok {
			return
		}
		// Now you can compare the UserID from the expense with the user's UUID
		if userUUID != expense.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not allowed to update this expense"})
			return
		}

		if err := c.ShouldBindJSON(&expense); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := db.Where("id = ?", id).Save(&expense).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update expense"})
			return
		}

		c.JSON(200, expense)
	}
}

func DeleteExpense(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id") // This is the expense's ID

		userUUID, expense, ok := getUserIDAndExpense(c, db, id)
		if !ok {
			return
		}
		// Now you can compare the UserID from the expense with the user's UUID
		if userUUID != expense.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not allowed to update this expense"})
			return
		}

		// If the user IDs match, delete the expense
		if err := db.Delete(&expense).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete expense"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
	}
}

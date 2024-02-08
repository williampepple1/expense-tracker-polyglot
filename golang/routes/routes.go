package routes

import (
	"expense-tracker/handlers"
	"expense-tracker/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {

	r.POST("/register", handlers.RegisterUser(db))
	r.POST("/login", handlers.LoginUser(db))

	// Use the authorization middleware for the following routes
	authorized := r.Group("/")
	authorized.Use(middleware.Authorize())
	{
		authorized.GET("/expenses", handlers.ListExpenses(db))
		authorized.GET("/expenses/user/:userId", handlers.ListUserExpenses(db))
		authorized.POST("/expenses", handlers.CreateExpense(db))
		authorized.GET("/expenses/:id", handlers.GetExpense(db))
		authorized.PUT("/expenses/:id", handlers.UpdateExpense(db))
		authorized.DELETE("/expenses/:id", handlers.DeleteExpense(db))
	}
}

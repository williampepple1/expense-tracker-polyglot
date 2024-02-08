package routes

import (
	"expense-tracker/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {

	r.POST("/register", handlers.RegisterUser(db))
	r.POST("/login", handlers.LoginUser(db))
}

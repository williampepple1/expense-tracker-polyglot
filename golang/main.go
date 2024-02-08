package main

import (
	"expense-tracker/config"
	"expense-tracker/models"
	"expense-tracker/routes"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const colorReset = "\033[0m"

var colors = []string{
	"\033[31m", // Red
	"\033[32m", // Green
	"\033[33m", // Yellow
	"\033[34m", // Blue
	"\033[35m", // Purple
	"\033[36m", // Cyan
}

func main() {
	port := "8080"
	env := "development"
	r := gin.Default()
	db, err := config.InitDB() // Initialize the database connection
	if err != nil {
		panic("Failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{})

	// Set up routes
	routes.SetupRoutes(r, db)

	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	color := colors[rng.Intn(len(colors))]

	r.Run()

	log.Println(color + "=================================" + colorReset)
	log.Printf("%s======= ENV: %s =======%s", color, env, colorReset)
	log.Printf("%s🚀 App listening on the port %s%s", color, port, colorReset)
	log.Println(color + "=================================" + colorReset)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}

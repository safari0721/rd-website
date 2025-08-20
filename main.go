package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/parthvinchhi/rd-website/pkg/db"
	"github.com/parthvinchhi/rd-website/pkg/handlers"
	"github.com/parthvinchhi/rd-website/pkg/models"
	"github.com/parthvinchhi/rd-website/pkg/repo"
	"github.com/parthvinchhi/rd-website/pkg/routes"
	"github.com/parthvinchhi/rd-website/pkg/services"
)

func main() {
	db, err := db.Connect()
	if err != nil {
		log.Fatalln("Database connection failed!")
	}

	log.Println("Connected to database successfully")

	if os.Getenv("DISABLE_AUTOMIGRATE") == "" {
		if err := db.AutoMigrate(&models.User{}, &models.Accounts{}); err != nil {
			log.Fatalln("Failed to migrate models")
		}
		log.Println("Migrated to database successfully")
	}

	userRepo := &repo.UserRepo{DB: db}
	authService := &services.AuthService{Repo: userRepo}
	authHandler := &handlers.AuthHandler{AuthService: authService}

	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	app := routes.SetupRouter(authHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Run(":" + port)
}

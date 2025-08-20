package main

import (
	"log"

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

	if err := db.AutoMigrate(&models.User{}, &models.Accounts{}); err != nil {
		log.Fatalln("Failed to migrate models")
	}

	log.Println("Migrated to database successfully")

	userRepo := &repo.UserRepo{DB: db}
	authService := &services.AuthService{Repo: userRepo}
	authHandler := &handlers.AuthHandler{AuthService: authService}

	routes := routes.SetupRouter(authHandler)
	routes.Run(":8080")
}

package main

import (
	"bookmarket/config"
	"bookmarket/models"
	"bookmarket/routes"

	"github.com/joho/godotenv"
)

func init() {
	// load .env file
	godotenv.Load(".env")
	db := config.SetupDB()

	// Migration table user
	db.AutoMigrate(&models.Users{})
}

func main() {
	// Call router
	router := routes.SetupRouter()

	router.Start(":8080")
}

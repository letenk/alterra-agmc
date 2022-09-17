package main

import (
	"bookmarket/config"
	"bookmarket/models"
	"bookmarket/routes"

	"github.com/joho/godotenv"
)

func main() {
	// load .env file
	godotenv.Load(".env")
	db := config.SetupDB()

	// Migration table user
	db.AutoMigrate(&models.Users{})

	// Call router
	router := routes.SetupRouter(db)

	router.Start(":8080")
}

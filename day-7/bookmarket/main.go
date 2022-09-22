package main

import (
	"bookmarket/database"
	"bookmarket/internal/app/users/models"
	"bookmarket/internal/http"
	m "bookmarket/internal/middlewares"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// load .env file
	godotenv.Load(".env")
	db := database.SetupDB()

	// Migration table user
	db.AutoMigrate(&models.Users{})

	// Create new router echo
	e := echo.New()
	// Use middleware
	m.LogMiddleware(e)
	// Use http
	http.NewHTTP(e, db)

	e.Logger.Fatal(e.Start(":8080"))
}

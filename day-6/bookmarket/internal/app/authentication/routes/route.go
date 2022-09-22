package routes

import (
	"bookmarket/internal/app/authentication/handlers"
	"bookmarket/internal/app/users/repository"
	"bookmarket/internal/app/users/service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Router(db *gorm.DB, g *echo.Group) {
	// Use repository
	userRepository := repository.NewRepository(db)
	// Use service
	userService := service.NewService(userRepository)
	// Use handler
	authenticationHandlers := handlers.NewHandlers(userService)

	// Create user
	g.POST("/login", authenticationHandlers.Login)
}

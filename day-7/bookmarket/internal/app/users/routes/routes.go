package routes

import (
	"bookmarket/internal/app/users/handlers"
	"bookmarket/internal/app/users/repository"
	"bookmarket/internal/app/users/service"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func Router(db *gorm.DB, g *echo.Group) {
	// Use repository
	userRepository := repository.NewRepository(db)
	// Use service
	userService := service.NewService(userRepository)
	// Use handler
	userhandlers := handlers.NewHandlers(userService)

	// Use middleware for users group prefix
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		AuthScheme: "Bearer",
		Skipper: func(c echo.Context) bool {
			// Skip middleware if method is equal 'POST'
			if c.Request().Method == "POST" {
				return true
			}
			return false
		},
	}))

	// Find All
	g.GET("", userhandlers.GetAll)
	// Find by id
	g.GET("/:id", userhandlers.FindByID)
	// Create user
	g.POST("", userhandlers.Create)
	// Update
	g.PUT("/:id", userhandlers.Update)
	// Delete
	g.DELETE("/:id", userhandlers.Delete)
}

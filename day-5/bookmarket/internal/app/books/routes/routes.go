package routes

import (
	"bookmarket/internal/app/books/handlers"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Router(g *echo.Group) {
	// Use handler
	bookHandlers := handlers.NewHandlers()

	// Use middleware for books group prefix
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		AuthScheme: "Bearer",
		Skipper: func(c echo.Context) bool {
			// Skip middleware if method is equal 'GET'
			if c.Request().Method == "GET" {
				return true
			}
			return false
		},
	}))

	// Endpoint get all books
	g.GET("", bookHandlers.GetAll)
	// Endpoint create new book
	g.POST("", bookHandlers.Create)
	// Endpoint get book by id
	g.GET("/:id", bookHandlers.FindById)
	// Endpoint update by id
	g.PUT("/:id", bookHandlers.Update)
	// Endpoint get delete by id
	g.DELETE("/:id", bookHandlers.Delete)

}

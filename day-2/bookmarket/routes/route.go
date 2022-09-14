package routes

import (
	"bookmarket/config"
	"bookmarket/controllers"

	"github.com/labstack/echo/v4"
)

func SetupRouter() *echo.Echo {
	// Create new instance
	router := echo.New()

	// Use book controller
	bookController := controllers.NewBookController()

	// Open connection db
	db := config.SetupDB()
	// Use user controller with argument var db
	userController := controllers.NewUserController(db)

	// Version group
	v1 := router.Group("/v1")

	// Books endpoint group
	books := v1.Group("/books")
	// Endpoint get all books
	books.GET("", bookController.GetAll)
	// Endpoint create new book
	books.POST("", bookController.Create)
	// Endpoint get book by id
	books.GET("/:id", bookController.FindById)
	// Endpoint update by id
	books.PUT("/:id", bookController.Update)
	// Endpoint get delete by id
	books.DELETE("/:id", bookController.Delete)

	// Users endpoint group
	users := v1.Group("/users")
	// Endpoint get all users
	users.GET("", userController.GetAll)
	// Endpoint create new user
	users.POST("", userController.Create)
	// Endpoint get user by id
	users.GET("/:id", userController.FindByID)
	// Endpoint update user by id
	users.PUT("/:id", userController.Update)
	// Endpoint delete user by id
	users.DELETE("/:id", userController.Delete)

	return router
}

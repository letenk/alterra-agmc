package routes

import (
	"bookmarket/config"
	"bookmarket/controllers"
	"bookmarket/lib"
	m "bookmarket/middlewares"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Iteration error
		errors := lib.ValidationError(err)
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, errors)
	}
	return nil
}

func SetupRouter() *echo.Echo {
	// Create new instance
	router := echo.New()

	// Use middleware log
	m.LogMiddleware(router)

	// Use validator
	router.Validator = &CustomValidator{validator: validator.New()}

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

	// Use middleware for books group prefix
	books.Use(middleware.JWTWithConfig(middleware.JWTConfig{
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

	// Use middleware for users group prefix
	users.Use(middleware.JWTWithConfig(middleware.JWTConfig{
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

	// Login endpoint
	v1.POST("/login", userController.Login)

	return router
}

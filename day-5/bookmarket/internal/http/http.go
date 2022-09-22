package http

import (
	authentication "bookmarket/internal/app/authentication/routes"
	books "bookmarket/internal/app/books/routes"
	users "bookmarket/internal/app/users/routes"
	"bookmarket/pkg/util"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Iteration error
		errors := util.ValidationError(err)
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, errors)
	}
	return nil
}

func NewHTTP(e *echo.Echo, db *gorm.DB) {
	// Validation
	e.Validator = &CustomValidator{validator: validator.New()}
	v1 := e.Group("/v1")
	authentication.Router(db, v1)
	users.Router(db, v1.Group("/users"))
	books.Router(v1.Group("/books"))
}

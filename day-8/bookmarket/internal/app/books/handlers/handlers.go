package handlers

import (
	"bookmarket/internal/app/books/models"
	"bookmarket/internal/dto"
	"bookmarket/pkg/util"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct{}

func NewHandlers() *handler {
	return &handler{}
}

// Create new var books with type slice based on models Books to hold data static
var books []models.Books

// Run function init for append sample data into var `books` on above
func init() {
	book := models.Books{
		ID:        "6d55b8f0-df37-4c38-9e5b-e780bba68381",
		Name:      "Automic Habbits",
		Author:    "James Clear",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	books = append(books, book)
}

func (h *handler) GetAll(ctx echo.Context) error {
	// Check length / contains from slice books
	if len(books) < 1 {
		// Create format api response for error
		response := util.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"Books not available, please insert first data",
			books,
		)

		return ctx.JSON(http.StatusBadRequest, response)
	}

	// Create format api response for success with return all data books
	response := util.ApiResponseWithData(
		http.StatusOK,
		"success",
		"List of books",
		books,
	)

	return ctx.JSON(http.StatusOK, response)
}

func (h *handler) FindById(ctx echo.Context) error {
	// Get params id
	bookID := ctx.Param("id")

	for index, item := range books {
		if item.ID == bookID {
			// Create format api response for success with return selected data books by id
			response := util.ApiResponseWithData(
				http.StatusOK,
				"success",
				"List of books",
				books[index],
			)
			return ctx.JSON(http.StatusOK, response)
		}
	}

	// Create format api response for error
	message := fmt.Sprintf("Book with id: %s not found", bookID)
	response := util.ApiResponseWithData(
		http.StatusBadRequest,
		"error",
		message,
		nil,
	)

	return ctx.JSON(http.StatusBadRequest, response)
}

func (h *handler) Create(ctx echo.Context) error {
	// newBook variable based on library CreateOrUpdateBook
	var payload dto.CreateOrUpdateBook

	// Binding payload request into var input
	err := ctx.Bind(&payload)
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := util.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"Create book failed",
			errors,
		)

		return ctx.JSON(http.StatusBadRequest, response)
	}

	// Validate
	err = ctx.Validate(payload)
	if err != nil {
		response := util.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"Create user failed",
			err,
		)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	var book models.Books
	// Generate uuid for id and binding into id book
	id := uuid.NewString()
	book.ID = id
	book.Name = payload.Name
	book.Author = payload.Author

	// Get time now and binding into field createdAt and updatedAt
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	// Append data newBook into var books
	books = append(books, book)

	response := util.ApiResponseWithData(
		http.StatusCreated,
		"success",
		"Book has been created",
		book,
	)

	return ctx.JSON(http.StatusCreated, response)
}

func (h *handler) Update(ctx echo.Context) error {
	// payload variable based on model books
	var payload dto.CreateOrUpdateBook

	// Binding payload request into var updatedBook
	err := ctx.Bind(&payload)
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := util.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"Update book failed",
			errors,
		)

		return ctx.JSON(http.StatusBadRequest, response)
	}

	// Validate
	err = ctx.Validate(payload)
	if err != nil {
		response := util.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"Update book failed",
			err,
		)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// Get params id
	bookID := ctx.Param("id")

	var book models.Books

	for index, item := range books {
		if item.ID == bookID {
			// Delete old book
			books = append(books[:index], books[index+1:]...)

			// Update or Re-Insert data book updated
			book.ID = item.ID
			book.Name = payload.Name
			book.Author = payload.Author
			book.CreatedAt = item.CreatedAt
			book.UpdatedAt = time.Now()
			books = append(books, book)

			// Create format api response for success with return selected data books by id
			response := util.ApiResponseWithData(
				http.StatusOK,
				"success",
				"Book has been updated",
				book,
			)
			return ctx.JSON(http.StatusOK, response)
		}
	}

	// Create format api response for error
	message := fmt.Sprintf("Book with id: %s not found", bookID)
	response := util.ApiResponseWithData(
		http.StatusBadRequest,
		"error",
		message,
		nil,
	)

	return ctx.JSON(http.StatusBadRequest, response)
}

func (h *handler) Delete(ctx echo.Context) error {
	// Get params id
	bookID := ctx.Param("id")

	for index, item := range books {
		if item.ID == bookID {
			// Delete old book
			books = append(books[:index], books[index+1:]...)

			// Create format api response for success with return selected data books by id
			response := util.ApiResponseWithData(
				http.StatusOK,
				"success",
				"Book has been deleted",
				nil,
			)
			return ctx.JSON(http.StatusOK, response)
		}
	}

	// Create format api response for error
	message := fmt.Sprintf("Book with id: %s not found", bookID)
	response := util.ApiResponseWithData(
		http.StatusBadRequest,
		"error",
		message,
		nil,
	)

	return ctx.JSON(http.StatusBadRequest, response)
}

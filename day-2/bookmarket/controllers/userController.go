package controllers

import (
	"bookmarket/lib"
	"bookmarket/models"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *userController {
	return &userController{db}
}

func (u *userController) GetAll(ctx echo.Context) error {
	var users []models.Users

	// Find all users and passing into var `users`
	err := u.db.Find(&users).Error
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := lib.ApiResponseWithData(
			http.StatusInternalServerError,
			"error",
			"internal server error",
			errors,
		)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	// Using formatter, for not displaying the password to response
	formatter := lib.FormatUsers(users)

	// If no error, create response and return JSON with data users
	response := lib.ApiResponseWithData(
		http.StatusOK,
		"success",
		"list of users",
		formatter,
	)

	return ctx.JSON(http.StatusOK, response)
}

func (u *userController) FindByID(ctx echo.Context) error {
	// Get parameters id
	id := ctx.Param("id")

	var user models.Users

	// Find by id
	err := u.db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := lib.ApiResponseWithData(
			http.StatusInternalServerError,
			"error",
			"internal server error",
			errors,
		)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	// If user not found
	if user.ID == "" {
		message := fmt.Sprintf("user with id: %s not found", id)
		// If no error, create response and return JSON with data users
		response := lib.ApiResponseWithData(
			http.StatusOK,
			"success",
			message,
			nil,
		)

		return ctx.JSON(http.StatusOK, response)
	}

	// Using formatter, for not displaying the password to response
	formatter := lib.FormatUser(user)

	// If no error, create response and return JSON with data users
	response := lib.ApiResponseWithData(
		http.StatusOK,
		"success",
		"Data of user",
		formatter,
	)

	return ctx.JSON(http.StatusOK, response)
}

func (u *userController) Create(ctx echo.Context) error {
	// newUser variable based on model users
	var newUser models.Users

	// Binding payload request into var newUser
	err := ctx.Bind(&newUser)
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := lib.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"Create user failed",
			errors,
		)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// Generate uuid for id and binding into id book
	id := uuid.NewString()
	newUser.ID = id

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := lib.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"Create user failed",
			errors,
		)
		return ctx.JSON(http.StatusBadRequest, response)
	}
	// passing password hash into object domain user
	newUser.Password = string(passwordHash)

	// Save user into db
	err = u.db.Save(&newUser).Error
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := lib.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"Create user failed",
			errors,
		)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// If no error, create response success and return JSON
	response := lib.ApiResponseWithData(
		http.StatusCreated,
		"success",
		"User has been created",
		newUser,
	)

	return ctx.JSON(http.StatusCreated, response)
}

func (u *userController) Update(ctx echo.Context) error {
	// Get parameters id
	id := ctx.Param("id")

	var updateUser models.Users

	// Binding payload request into var updateUser
	err := ctx.Bind(&updateUser)
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := lib.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"Update user failed",
			errors,
		)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	var user models.Users
	// Find by id
	err = u.db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := lib.ApiResponseWithData(
			http.StatusInternalServerError,
			"error",
			"internal server error",
			errors,
		)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	// If user not found
	if user.ID == "" {
		message := fmt.Sprintf("user with id: %s not found", id)
		// If no error, create response and return JSON with data users
		response := lib.ApiResponseWithData(
			http.StatusBadRequest,
			"success",
			message,
			nil,
		)

		return ctx.JSON(http.StatusBadRequest, response)
	}

	// If user is available, update data from payload request
	user.Fullname = updateUser.Fullname
	// If payload updatedUser not empty string, update the password
	if updateUser.Password != "" {
		// Hash password
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), bcrypt.DefaultCost)
		if err != nil {
			errors := map[string]any{
				"errors": err.Error(),
			}
			response := lib.ApiResponseWithData(
				http.StatusBadRequest,
				"error",
				"updated user failed",
				errors,
			)
			return ctx.JSON(http.StatusBadRequest, response)
		}
		// passing password hash into object domain user
		user.Password = string(passwordHash)
	}
	// Update time propert updatedAt
	user.UpdatedAt = time.Now()

	// Save again with new data update
	err = u.db.Save(&user).Error
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := lib.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"updated user failed",
			errors,
		)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// Using formatter, for not displaying the password to response
	formatter := lib.FormatUser(user)

	// If no error, create response success and return JSON
	response := lib.ApiResponseWithData(
		http.StatusOK,
		"success",
		"User has been updated",
		formatter,
	)

	return ctx.JSON(http.StatusOK, response)
}

func (u *userController) Delete(ctx echo.Context) error {
	// Get parameters id
	id := ctx.Param("id")

	var user models.Users
	// Find by id
	err := u.db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := lib.ApiResponseWithData(
			http.StatusInternalServerError,
			"error",
			"internal server error",
			errors,
		)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	// If user not found
	if user.ID == "" {
		message := fmt.Sprintf("user with id: %s not found", id)
		// If no error, create response and return JSON with data users
		response := lib.ApiResponseWithData(
			http.StatusBadRequest,
			"success",
			message,
			nil,
		)

		return ctx.JSON(http.StatusBadRequest, response)
	}

	// delete user
	err = u.db.Delete(&user).Error
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := lib.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"delete user failed",
			errors,
		)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// If no error, create response success and return JSON
	response := lib.ApiResponseWithData(
		http.StatusOK,
		"success",
		"User has been deleted",
		nil,
	)

	return ctx.JSON(http.StatusOK, response)
}

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
		lib.HandleInternalServerError(ctx, err)
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
		lib.HandleInternalServerError(ctx, err)
	}

	// If user not found
	if user.ID == "" {
		message := fmt.Sprintf("user with id: %s not found", id)
		// If no error, create response and return JSON with data users
		response := lib.ApiResponseWithData(
			http.StatusOK,
			"error",
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
	// input variable based on formatter CreateOrUpdateUser
	var input lib.CreateOrUpdateUser

	// Binding payload request into var newUser
	err := ctx.Bind(&input)
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

	// Validate
	err = ctx.Validate(input)
	if err != nil {
		response := lib.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"Create user failed",
			err,
		)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// Find user by email, for check whether user is available on the database
	var user models.Users
	err = u.db.Where("email = ?", input.Email).Find(&user).Error
	if err != nil {
		lib.HandleInternalServerError(ctx, err)
	}

	// If user.ID not empty string / user is available, return response error
	if user.ID != "" {
		errors := map[string]any{
			"errors": "email already exist",
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
	user.ID = id
	user.Email = input.Email
	user.Fullname = input.Fullname
	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
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
	user.Password = string(passwordHash)

	// Save user into db
	err = u.db.Save(&user).Error
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
		user,
	)

	return ctx.JSON(http.StatusCreated, response)
}

func (u *userController) Update(ctx echo.Context) error {
	// Get token from header `Authorization`
	token := ctx.Request().Header.Get("Authorization")

	// Parse Token and get only current id user is logged in
	currentIdUser, err := lib.ParseTokenJWT(token)
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := lib.ApiResponseWithData(
			http.StatusUnauthorized,
			"error",
			"unauthorized",
			errors,
		)
		return ctx.JSON(http.StatusUnauthorized, response)
	}

	// Get parameters id
	id := ctx.Param("id")

	// Check whether `current user logged in` is not same with params `id`
	if currentIdUser != id {
		errors := map[string]any{
			"errors": "not access",
		}
		response := lib.ApiResponseWithData(
			http.StatusUnauthorized,
			"error",
			"unauthorized",
			errors,
		)
		return ctx.JSON(http.StatusUnauthorized, response)
	}

	var input lib.CreateOrUpdateUser

	// Binding payload request into var updateUser
	err = ctx.Bind(&input)
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
		lib.HandleInternalServerError(ctx, err)
	}

	// If user not found
	if user.ID == "" {
		message := fmt.Sprintf("user with id: %s not found", id)
		// If no error, create response and return JSON with data users
		response := lib.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			message,
			nil,
		)

		return ctx.JSON(http.StatusBadRequest, response)
	}

	// If user is available, update data from payload request
	user.Fullname = input.Fullname
	// If payload updatedUser not empty string, update the password
	if input.Password != "" {
		// Hash password
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
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
	// Get token from header `Authorization`
	token := ctx.Request().Header.Get("Authorization")

	// Parse Token and get only current id user is logged in
	currentIdUser, err := lib.ParseTokenJWT(token)
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := lib.ApiResponseWithData(
			http.StatusUnauthorized,
			"error",
			"unauthorized",
			errors,
		)
		return ctx.JSON(http.StatusUnauthorized, response)
	}

	// Get parameters id
	id := ctx.Param("id")

	// Check whether `current user logged in` is not same with params `id`
	if currentIdUser != id {
		errors := map[string]any{
			"errors": "not access",
		}
		response := lib.ApiResponseWithData(
			http.StatusUnauthorized,
			"error",
			"unauthorized",
			errors,
		)
		return ctx.JSON(http.StatusUnauthorized, response)
	}

	var user models.Users
	// Find by id
	err = u.db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		lib.HandleInternalServerError(ctx, err)
	}

	// If user not found
	if user.ID == "" {
		message := fmt.Sprintf("user with id: %s not found", id)
		// If no error, create response and return JSON with data users
		response := lib.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
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

func (u *userController) Login(ctx echo.Context) error {
	var userLogin models.Users
	// Binding payload request into var userLogin
	err := ctx.Bind(&userLogin)
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := lib.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"bad request",
			errors,
		)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// Find user by email
	var user models.Users
	err = u.db.Where("email = ?", userLogin.Email).Find(&user).Error
	if err != nil {
		lib.HandleInternalServerError(ctx, err)
	}

	// If user not found
	if user.ID == "" {
		// If no error, create response and return JSON with data users
		response := lib.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"user or password incorrect",
			nil,
		)

		return ctx.JSON(http.StatusBadRequest, response)
	}

	// If user is found, Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
	if err != nil {
		lib.HandleInternalServerError(ctx, err)
	}

	// Generate JWT
	token := lib.GenerateToken(ctx, user)

	// If no error, create response success and return JSON
	formatter := map[string]any{"token": token}
	response := lib.ApiResponseWithData(
		http.StatusOK,
		"success",
		"You are logged",
		formatter,
	)

	return ctx.JSON(http.StatusOK, response)
}

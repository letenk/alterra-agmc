package handlers

import (
	"bookmarket/internal/app/users/service"
	"bookmarket/internal/dto"
	"bookmarket/pkg/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service service.Service
}

func NewHandlers(service service.Service) *handler {
	return &handler{service}
}

func (h *handler) GetAll(ctx echo.Context) error {
	// Get all user
	users, err := h.service.GetAll(ctx.Request().Context())
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := util.ApiResponseWithData(
			http.StatusInternalServerError,
			"error",
			"internal server error",
			errors,
		)

		return ctx.JSON(http.StatusInternalServerError, response)
	}

	// Using formatter, for not displaying the password to response
	formatUserResponse := dto.FormatUsersResponse(users)

	// If no error, create response and return JSON with data users
	response := util.ApiResponseWithData(
		http.StatusOK,
		"success",
		"list of users",
		formatUserResponse,
	)

	return ctx.JSON(http.StatusOK, response)
}

func (h *handler) FindByID(ctx echo.Context) error {
	// Get parameters id
	id := ctx.Param("id")

	// Find user by id
	user, err := h.service.FindByID(ctx.Request().Context(), id)
	if user.ID == "" {
		if err != nil {
			message := fmt.Sprintf("user with id: %s not found", id)
			errors := map[string]any{
				"errors": message,
			}
			response := util.ApiResponseWithData(
				http.StatusInternalServerError,
				"error",
				"internal server error",
				errors,
			)

			return ctx.JSON(http.StatusInternalServerError, response)
		}
	}

	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := util.ApiResponseWithData(
			http.StatusInternalServerError,
			"error",
			"internal server error",
			errors,
		)

		return ctx.JSON(http.StatusInternalServerError, response)
	}

	// Using formatter, for not displaying the password to response
	formatUserResponse := dto.FormatUserResponse(user)

	// If no error, create response and return JSON with data users
	response := util.ApiResponseWithData(
		http.StatusOK,
		"success",
		"list of users",
		formatUserResponse,
	)

	return ctx.JSON(http.StatusOK, response)
}

func (h *handler) Create(ctx echo.Context) error {
	var payload dto.CreateRequestBody

	// Binding payload request into var newUser
	err := ctx.Bind(&payload)
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := util.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"Create user failed",
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

	// Find user by email, for check whether user is available on the database
	user, _ := h.service.FindByEmail(ctx.Request().Context(), payload.Email)
	if user.ID != "" {
		errors := map[string]any{
			"errors": "email already exist",
		}
		response := util.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"Create user failed",
			errors,
		)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// If email is not available, create new user
	newUser, err := h.service.Create(ctx.Request().Context(), &payload)
	if err != nil {
		response := util.ApiResponseWithData(
			http.StatusInternalServerError,
			"internal server error",
			"Create user failed",
			nil,
		)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// Using formatter, for not displaying the password to response
	formatUserResponse := dto.FormatUserResponse(newUser)

	// If no error, create response success and return JSON
	response := util.ApiResponseWithData(
		http.StatusCreated,
		"success",
		"User has been created",
		formatUserResponse,
	)

	return ctx.JSON(http.StatusCreated, response)
}

func (h *handler) Update(ctx echo.Context) error {
	// Get token from header `Authorization`
	token := ctx.Request().Header.Get("Authorization")

	// Parse Token and get only current id user is logged in
	currentIdUser, err := util.ParseTokenJWT(token)
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := util.ApiResponseWithData(
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
		response := util.ApiResponseWithData(
			http.StatusUnauthorized,
			"error",
			"unauthorized",
			errors,
		)
		return ctx.JSON(http.StatusUnauthorized, response)
	}

	var payload dto.UpdateRequestBody
	// Binding payload request into var updateUser
	err = ctx.Bind(&payload)
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := util.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"Update user failed",
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

	// Find user by id
	userUpdated, err := h.service.Update(ctx.Request().Context(), id, &payload)
	if err != nil {
		// If no error, create response and return JSON with data users
		response := util.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			err.Error(),
			nil,
		)

		return ctx.JSON(http.StatusBadRequest, response)
	}

	// Using formatter, for not displaying the password to response
	formatter := dto.FormatUserResponse(userUpdated)

	// If no error, create response success and return JSON
	response := util.ApiResponseWithData(
		http.StatusOK,
		"success",
		"User has been updated",
		formatter,
	)

	return ctx.JSON(http.StatusOK, response)
}

func (h *handler) Delete(ctx echo.Context) error {
	// Get token from header `Authorization`
	token := ctx.Request().Header.Get("Authorization")

	// Parse Token and get only current id user is logged in
	currentIdUser, err := util.ParseTokenJWT(token)
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := util.ApiResponseWithData(
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
		response := util.ApiResponseWithData(
			http.StatusUnauthorized,
			"error",
			"unauthorized",
			errors,
		)
		return ctx.JSON(http.StatusUnauthorized, response)
	}

	// Delete
	err = h.service.Delete(ctx.Request().Context(), id)
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := util.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"delete user failed",
			errors,
		)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// If no error, create response success and return JSON
	response := util.ApiResponseWithData(
		http.StatusOK,
		"success",
		"User has been deleted",
		nil,
	)

	return ctx.JSON(http.StatusOK, response)
}

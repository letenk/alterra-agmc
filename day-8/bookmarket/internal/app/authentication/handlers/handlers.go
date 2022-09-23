package handlers

import (
	"bookmarket/internal/app/users/service"
	"bookmarket/internal/dto"
	"bookmarket/pkg/util"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type handler struct {
	service service.Service
}

func NewHandlers(service service.Service) *handler {
	return &handler{service}
}

func (h *handler) Login(ctx echo.Context) error {
	var payload dto.LoginPayload
	// Binding payload request into var userLogin
	err := ctx.Bind(&payload)
	if err != nil {
		errors := map[string]any{
			"errors": err.Error(),
		}
		response := util.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"bad request",
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
			"bad request",
			err,
		)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// Find user by email
	user, _ := h.service.FindByEmail(ctx.Request().Context(), payload.Email)

	// If user not found
	if user.ID == "" {
		// If no error, create response and return JSON with data users
		response := util.ApiResponseWithData(
			http.StatusBadRequest,
			"error",
			"user or password incorrect",
			nil,
		)

		return ctx.JSON(http.StatusBadRequest, response)
	}

	// If user is found, Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		response := util.ApiResponseWithData(
			http.StatusInternalServerError,
			"error",
			"internal server error",
			err,
		)

		return ctx.JSON(http.StatusInternalServerError, response)
	}

	// Generate JWT
	token, err := util.GenerateToken(ctx, user)
	if err != nil {
		response := util.ApiResponseWithData(
			http.StatusInternalServerError,
			"error",
			"internal server error",
			err,
		)

		return ctx.JSON(http.StatusInternalServerError, response)
	}

	// If no error, create response success and return JSON
	formatter := map[string]any{"token": token}
	response := util.ApiResponseWithData(
		http.StatusOK,
		"success",
		"You are logged",
		formatter,
	)

	return ctx.JSON(http.StatusOK, response)
}

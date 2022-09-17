package lib

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleInternalServerError(ctx echo.Context, err error) any {
	errors := map[string]any{
		"errors": err.Error(),
	}
	response := ApiResponseWithData(
		http.StatusInternalServerError,
		"error",
		"internal server error",
		errors,
	)

	return ctx.JSON(http.StatusInternalServerError, response)
}

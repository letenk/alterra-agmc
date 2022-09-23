package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func LogMiddleware(e *echo.Echo) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		// Format log
		Format: `time=${time_rfc3339} method=${method} path=${host}${uri} status=${status} latency=${latency_human} error=${error}` + "\n",
	}))
}

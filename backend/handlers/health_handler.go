package handlers

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

// health check to see if the server is running
func HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "healthy",
		"message": "Server is running fine .....",
	})
}
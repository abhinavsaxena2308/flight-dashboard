package handlers

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

// HealthHandler handles the health check endpoint
func HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "healthy",
		"message": "Server is running",
	})
}
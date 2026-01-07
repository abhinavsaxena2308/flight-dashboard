package routes

import (
	"flight-dashboard-backend/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
)

// SetupRoutes configures all API routes for the flight dashboard backend
// It registers endpoints for health checks, state-wise flight data, and other APIs
func SetupRoutes(e *echo.Echo) {
	// Health check endpoint
	e.GET("/health", handlers.HealthHandler)

	// Default endpoint
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to Flight Dashboard Backend API",
		})
	})

	// State-wise flight data endpoints
	e.GET("/api/state-flights", handlers.GetStateWiseFlights)
	e.GET("/api/states", handlers.GetStateList)
	e.GET("/api/state/:state", handlers.GetStateDetail)
	e.GET("/api/states/:state/airlines", handlers.GetTopAirlinesForState)
}

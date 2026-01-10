package routes

import (
	"flight-dashboard-backend/handlers"
	"net/http"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/health", handlers.HealthHandler)     // health check endpoint
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{		//default endpoint
			"message": "Welcome to Flight Dashboard Backend API",
		})
	})

	// state-wise flight data endpoints
	e.GET("/api/state-flights", handlers.GetStateWiseFlights)
	e.GET("/api/states", handlers.GetStateList)
	e.GET("/api/state/:state", handlers.GetStateDetail)
	e.GET("/api/states/:state/airlines", handlers.GetTopAirlinesForState)
}

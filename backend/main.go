package main

import (
	"flight-dashboard-backend/routes"
	"flight-dashboard-backend/services"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// main is the entry point for the flight dashboard backend service
// It initializes all services, loads data, and starts the HTTP server
func main() {
	// Initialize city state mapper
	services.GetCityStateMapper()
	log.Println("City-to-state mapping initialized")

	// Initialize flight data service and load CSV data
	dataService := services.GetFlightDataService()
	csvPath := "data/dataset.csv"
	err := dataService.LoadFlightDataFromCSV(csvPath)
	if err != nil {
		log.Printf("Warning: Could not load flight data: %v", err)
		log.Println("Please ensure the CSV file exists in the data directory")
	} else {
		log.Printf("Successfully loaded %d flight records", dataService.GetFlightCount())
	}

	// Initialize state aggregator (precomputes all state-wise aggregations)
	services.GetStateAggregator()
	log.Println("State-wise aggregations computed and stored in memory")

	// Configure Echo framework with middleware
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS middleware
	e.Use(middleware.CORS())

	// Setup routes
	routes.SetupRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

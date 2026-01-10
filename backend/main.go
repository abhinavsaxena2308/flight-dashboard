package main

import (
	"flight-dashboard-backend/routes"
	"flight-dashboard-backend/services"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)
func main() {
	services.GetCityStateMapper()
	//log.Println("City-to-state mapping initialized")

	// starting flight data service and load CSV data
	dataService := services.GetFlightDataService()
	csvPath := "data/dataset.csv"
	err := dataService.LoadFlightDataFromCSV(csvPath)
	if err != nil {
		log.Printf("Warning: Could not load flight data: %v", err)
		log.Println("Please ensure the CSV file exists in the data directory")
	} else {
		log.Printf("Successfully loaded %d flight records", dataService.GetFlightCount())
	}

	// starting state aggregator (precomputes all state-wise aggregations)
	services.GetStateAggregator()
	//log.Println("State-wise aggregations computed and stored in memory")
	e := echo.New()  //echo-fw
	e.Use(middleware.Logger())  //middleware
	e.Use(middleware.Recover())  //middleware
	e.Use(middleware.CORS())  //cors
	routes.SetupRoutes(e)  //routes
	e.Logger.Fatal(e.Start(":8080"))  //port-ini
}

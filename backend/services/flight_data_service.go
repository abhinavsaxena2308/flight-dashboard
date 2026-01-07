package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"flight-dashboard-backend/models"
)

// FlightDataService handles loading and accessing flight data
type FlightDataService struct {
	flights []models.Flight
	mutex   sync.RWMutex
}

// Global instance of the flight data service
var flightDataService *FlightDataService
var once sync.Once

// GetFlightDataService returns a singleton instance of FlightDataService
func GetFlightDataService() *FlightDataService {
	once.Do(func() {
		flightDataService = &FlightDataService{}
	})
	return flightDataService
}

// LoadFlightDataFromCSV loads flight data from a CSV file into memory
func (fds *FlightDataService) LoadFlightDataFromCSV(filePath string) error {
	fds.mutex.Lock()
	defer fds.mutex.Unlock()

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header to understand column order
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %v", err)
	}

	// Map header columns to their indices
	columnIndices := make(map[string]int)
	for i, col := range header {
		columnIndices[strings.TrimSpace(strings.ToLower(col))] = i
	}

	var flights []models.Flight

	// Read each record
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading CSV record: %v", err)
			continue // Skip invalid records
		}

		flight, err := parseFlightRecord(record, columnIndices)
		if err != nil {
			log.Printf("Skipping invalid record: %v, Record: %v", err, record)
			continue // Skip invalid records
		}

		flights = append(flights, flight)
	}

	fds.flights = flights
	log.Printf("Successfully loaded %d flight records from %s", len(flights), filePath)
	return nil
}

// parseFlightRecord converts a CSV record to a Flight struct
func parseFlightRecord(record []string, indices map[string]int) (models.Flight, error) {
	var flight models.Flight

	// Helper function to safely get field value
	getField := func(fieldName string) string {
		if idx, exists := indices[fieldName]; exists && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
		return ""
	}

	// Map fields based on common CSV column names
	flight.Airline = getField("airline")
	flight.FlightDate = getField("date_of_journey") // Common name in flight datasets
	if flight.FlightDate == "" {
		flight.FlightDate = getField("flight_date") // Alternative name
	}
	if flight.FlightDate == "" {
		flight.FlightDate = getField("flight date") // CSV column name
	}

	flight.Source = getField("source")
	if flight.Source == "" {
		flight.Source = getField("from_city") // Alternative name
	}
	if flight.Source == "" {
		flight.Source = getField("from") // CSV column name
	}

	flight.Destination = getField("destination")
	if flight.Destination == "" {
		flight.Destination = getField("to_city") // Alternative name
	}
	if flight.Destination == "" {
		flight.Destination = getField("to") // CSV column name
	}

	flight.FlightClass = getField("class")
	if flight.FlightClass == "" {
		flight.FlightClass = getField("flight_class") // Alternative name
	}

	// Parse duration (may be in hours or as a string)
	durationStr := getField("duration")
	if durationStr != "" {
		flight.Duration = parseDuration(durationStr)
	}

	// Parse price
	priceStr := getField("price")
	if priceStr != "" {
		// Remove commas and other non-numeric characters (except decimal point)
		cleanPriceStr := strings.ReplaceAll(priceStr, ",", "")
		// Remove any non-numeric characters except decimal point
		cleanPriceStr = cleanNonNumeric(cleanPriceStr)
		price, err := strconv.ParseFloat(cleanPriceStr, 64)
		if err != nil {
			price = 0 // Default to 0 if parsing fails
		}
		flight.Price = price
	}

	flight.DepartureTime = getField("departure_time")
	if flight.DepartureTime == "" {
		flight.DepartureTime = getField("dep_time") // Alternative name
	}
	if flight.DepartureTime == "" {
		flight.DepartureTime = getField("dep_time") // CSV column name (as 'from' field)
	}

	flight.ArrivalTime = getField("arrival_time")
	if flight.ArrivalTime == "" {
		flight.ArrivalTime = getField("arrival_time") // Alternative name
	}
	if flight.ArrivalTime == "" {
		flight.ArrivalTime = getField("arr_time") // CSV column name
	}

	// Parse stops
	stopsStr := getField("stops")
	if stopsStr != "" {
		// Clean the stops string to remove non-numeric characters
		cleanStopsStr := cleanNonNumeric(stopsStr)
		stops, err := strconv.Atoi(cleanStopsStr)
		if err != nil {
			stops = 0 // Default to 0 if parsing fails
		}
		flight.Stops = stops
	}

	flight.AdditionalInfo = getField("additional_info")
	if flight.AdditionalInfo == "" {
		flight.AdditionalInfo = getField("info") // Alternative name
	}

	return flight, nil
}

// GetAllFlights returns all loaded flight records
func (fds *FlightDataService) GetAllFlights() []models.Flight {
	fds.mutex.RLock()
	defer fds.mutex.RUnlock()

	// Return a copy to prevent external modification
	flightsCopy := make([]models.Flight, len(fds.flights))
	copy(flightsCopy, fds.flights)
	return flightsCopy
}

// GetFlightCount returns the total number of loaded flights
func (fds *FlightDataService) GetFlightCount() int {
	fds.mutex.RLock()
	defer fds.mutex.RUnlock()
	return len(fds.flights)
}

// GetStateForCity wraps the city state mapper to get state for a given city
func (fds *FlightDataService) GetStateForCity(city string) (string, bool) {
	mapper := GetCityStateMapper() // Use the new city state mapper
	return mapper.GetStateForCity(city)
}

// GetFlightCountByState returns the number of flights for a specific state
func (fds *FlightDataService) GetFlightCountByState(state string) int {
	count := 0
	state = strings.ToLower(state)

	fds.mutex.RLock()
	defer fds.mutex.RUnlock()

	for _, flight := range fds.flights {
		// Check if source or destination is in the given state
		if sourceState, ok := fds.GetStateForCity(flight.Source); ok && strings.ToLower(sourceState) == state {
			count++
		} else if destState, ok := fds.GetStateForCity(flight.Destination); ok && strings.ToLower(destState) == state {
			count++
		}
	}
	return count
}

// parseDuration parses duration strings like "2h 15m", "2h", "15m", "2.5h", etc.
func parseDuration(durationStr string) float64 {
	// Remove extra spaces
	durationStr = strings.TrimSpace(durationStr)

	// Handle special cases like "non-stop", "1-stop", etc.
	if strings.Contains(strings.ToLower(durationStr), "non-stop") {
		return 0 // Non-stop flights have no additional duration from stops
	}
	if strings.Contains(strings.ToLower(durationStr), "-stop") {
		// For values like "1-stop", "2-stop", etc., return 0 as duration
		return 0
	}

	// Handle "Xh Ym" format
	if strings.Contains(durationStr, " ") {
		parts := strings.Split(durationStr, " ")
		var totalMinutes float64 = 0

		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.Contains(part, "h") {
				hoursStr := strings.ReplaceAll(part, "h", "")
				hours, err := strconv.ParseFloat(strings.TrimSpace(hoursStr), 64)
				if err == nil {
					totalMinutes += hours * 60 // Convert to minutes for now
				}
			} else if strings.Contains(part, "m") {
				minutesStr := strings.ReplaceAll(part, "m", "")
				minutes, err := strconv.ParseFloat(strings.TrimSpace(minutesStr), 64)
				if err == nil {
					totalMinutes += minutes
				}
			}
		}
		return totalMinutes / 60 // Convert back to hours
	}

	// Handle "Xh" or "Ym" format
	if strings.Contains(durationStr, "h") {
		hoursStr := strings.ReplaceAll(durationStr, "h", "")
		hours, err := strconv.ParseFloat(strings.TrimSpace(hoursStr), 64)
		if err != nil {
			return 0
		}
		return hours
	}

	if strings.Contains(durationStr, "m") {
		minutesStr := strings.ReplaceAll(durationStr, "m", "")
		minutes, err := strconv.ParseFloat(strings.TrimSpace(minutesStr), 64)
		if err != nil {
			return 0
		}
		return minutes / 60 // Convert minutes to hours
	}

	// If it's just a number (already in hours), try to parse as float
	value, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0
	}
	return value
}

// cleanNonNumeric removes any non-numeric characters except decimal point
func cleanNonNumeric(s string) string {
	var cleaned strings.Builder
	for _, r := range s {
		if (r >= '0' && r <= '9') || r == '.' {
			cleaned.WriteRune(r)
		}
	}
	return cleaned.String()
}

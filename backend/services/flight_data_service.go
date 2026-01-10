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

// this service handles loading and accessing the flight data from CSV
type FlightDataService struct {
	flights []models.Flight
	mutex   sync.RWMutex
}

// global instance so we can access the flight data anywhere in the app
var flightDataService *FlightDataService
var once sync.Once

// returns a singleton instance of the flight data service - ensures we only have one instance
func GetFlightDataService() *FlightDataService {
	once.Do(func() {
		flightDataService = &FlightDataService{}
	})
	return flightDataService
}

// loads flight data from CSV file into memory - this is called when the app starts
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
			continue 
		}

		flight, err := parseFlightRecord(record, columnIndices)
		if err != nil {
			log.Printf("Skipping invalid record: %v, Record: %v", err, record)
			continue 
		}

		flights = append(flights, flight)
	}

	fds.flights = flights
	log.Printf("Successfully loaded %d flight records from %s", len(flights), filePath)
	return nil
}

// converts a CSV record to a Flight struct 
func parseFlightRecord(record []string, indices map[string]int) (models.Flight, error) {
	var flight models.Flight

	// Helper function to get field value
	getField := func(fieldName string) string {
		if idx, exists := indices[fieldName]; exists && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
		return ""
	}

	// Map fields 
	flight.Airline = getField("airline")
	flight.FlightDate = getField("date_of_journey") 
	if flight.FlightDate == "" {
		flight.FlightDate = getField("flight_date") 
	}
	if flight.FlightDate == "" {
		flight.FlightDate = getField("flight date") 
	}

	flight.Source = getField("source")
	if flight.Source == "" {
		flight.Source = getField("from_city") 
	}
	if flight.Source == "" {
		flight.Source = getField("from") 
	}

	flight.Destination = getField("destination")
	if flight.Destination == "" {
		flight.Destination = getField("to_city") 
	}
	if flight.Destination == "" {
		flight.Destination = getField("to") 
	}

	flight.FlightClass = getField("class")
	if flight.FlightClass == "" {
		flight.FlightClass = getField("flight_class") 
	}

	// flight duration
	durationStr := getField("duration")
	if durationStr != "" {
		flight.Duration = parseDuration(durationStr)
	}

	// Parse price
	priceStr := getField("price")
	if priceStr != "" {
		cleanPriceStr := strings.ReplaceAll(priceStr, ",", "")
		cleanPriceStr = cleanNonNumeric(cleanPriceStr)
		price, err := strconv.ParseFloat(cleanPriceStr, 64)
		if err != nil {
			price = 0 // default to 0 if parsing fails
		}
		flight.Price = price
	}

	flight.DepartureTime = getField("departure_time")
	if flight.DepartureTime == "" {
		flight.DepartureTime = getField("dep_time") 
	}
	if flight.DepartureTime == "" {
		flight.DepartureTime = getField("dep_time") 
	}

	flight.ArrivalTime = getField("arrival_time")
	if flight.ArrivalTime == "" {
		flight.ArrivalTime = getField("arrival_time") 
	}
	if flight.ArrivalTime == "" {
		flight.ArrivalTime = getField("arr_time") 
	}

	// parse stops.....
	stopsStr := getField("stops")
	if stopsStr != "" {
		cleanStopsStr := cleanNonNumeric(stopsStr)
		stops, err := strconv.Atoi(cleanStopsStr)
		if err != nil {
			stops = 0 // default to 0 if parsing fails
		}
		flight.Stops = stops
	}

	flight.AdditionalInfo = getField("additional_info")
	if flight.AdditionalInfo == "" {
		flight.AdditionalInfo = getField("info") 
	}

	return flight, nil
}

// returns all the loaded flight records
func (fds *FlightDataService) GetAllFlights() []models.Flight {
	fds.mutex.RLock()
	defer fds.mutex.RUnlock()

	// return a copy to prevent external modification of the original data
	flightsCopy := make([]models.Flight, len(fds.flights))
	copy(flightsCopy, fds.flights)
	return flightsCopy
}

// returns the total number of flights we've loaded - useful for stats
func (fds *FlightDataService) GetFlightCount() int {
	fds.mutex.RLock()
	defer fds.mutex.RUnlock()
	return len(fds.flights)
}

// gets the state for a given city using the city state mapper
func (fds *FlightDataService) GetStateForCity(city string) (string, bool) {
	mapper := GetCityStateMapper() // Use the new city state mapper
	return mapper.GetStateForCity(city)
}

// counts flights for a specific state - checks both source and destination
func (fds *FlightDataService) GetFlightCountByState(state string) int {
	count := 0
	state = strings.ToLower(state)

	fds.mutex.RLock()
	defer fds.mutex.RUnlock()

	for _, flight := range fds.flights {
		// checking if source or destination is in the given state
		if sourceState, ok := fds.GetStateForCity(flight.Source); ok && strings.ToLower(sourceState) == state {
			count++
		} else if destState, ok := fds.GetStateForCity(flight.Destination); ok && strings.ToLower(destState) == state {
			count++
		}
	}
	return count
}

// parses duration strings like '2h 15m', '2h', '15m', '2.5h', etc. - handles different formats
func parseDuration(durationStr string) float64 {
	// Remove extra spaces
	durationStr = strings.TrimSpace(durationStr)
	if strings.Contains(strings.ToLower(durationStr), "non-stop") {
		return 0 
	}
	if strings.Contains(strings.ToLower(durationStr), "-stop") {
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
					totalMinutes += hours * 60 
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

// removes non-numeric chars except decimal point - used for cleaning price/duration fields
func cleanNonNumeric(s string) string {
	var cleaned strings.Builder
	for _, r := range s {
		if (r >= '0' && r <= '9') || r == '.' {
			cleaned.WriteRune(r)
		}
	}
	return cleaned.String()
}

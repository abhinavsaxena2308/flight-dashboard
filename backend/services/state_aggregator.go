package services

import (
	"log"
	"strings"
	"sync"
)

type StateAggregation struct {
	StateName       string         `json:"state_name"`       // Name of the state
	TotalFlights    int            `json:"total_flights"`    // Total number of flights associated with this state
	IncomingFlights int            `json:"incoming_flights"` // Number of flights arriving to this state
	OutgoingFlights int            `json:"outgoing_flights"` // Number of flights departing from this state
	UniqueRoutes    int            `json:"unique_routes"`    // Number of unique routes in this state
	Airlines        map[string]int `json:"airlines"`         // Map of airline names to their flight counts
	RouteDetails    map[string]int `json:"route_details"`    // Map of route strings ("source->dest") to their counts
}

type StateAggregator struct {
	aggregations map[string]*StateAggregation // Map of state names to their aggregations
	mutex        sync.RWMutex                 // Mutex for thread-safe access to aggregations
	dataService  *FlightDataService           // Reference to flight data service
	mapper       *CityStateMapper             // Reference to city-to-state mapper
}

// Global instance of the state aggregator
var stateAggregator *StateAggregator
var aggOnce sync.Once

// GetStateAggregator returns a singleton instance of StateAggregator
func GetStateAggregator() *StateAggregator {
	aggOnce.Do(func() {
		stateAggregator = &StateAggregator{
			aggregations: make(map[string]*StateAggregation),
			dataService:  GetFlightDataService(),
			mapper:       GetCityStateMapper(),
		}
		// Precompute aggregations at startup
		stateAggregator.ComputeAggregations()
	})
	return stateAggregator
}

func (sa *StateAggregator) ComputeAggregations() {
	sa.mutex.Lock()
	defer sa.mutex.Unlock()

	// Get all flights
	flights := sa.dataService.GetAllFlights()

	// Initialize aggregation map
	aggregations := make(map[string]*StateAggregation)

	// Iterate through all flights to compute aggregations
	for _, flight := range flights {
		// Get states for source and destination
		// Normalize city names to match the mapper's expected format
		sourceState, sourceOk := sa.mapper.GetStateForCity(flight.Source)
		destState, destOk := sa.mapper.GetStateForCity(flight.Destination)

		// If not found with original name, try with normalized name
		if !sourceOk {
			sourceState, sourceOk = sa.mapper.GetStateForCity(normalizeCityNameForMapping(flight.Source))
		}
		if !destOk {
			destState, destOk = sa.mapper.GetStateForCity(normalizeCityNameForMapping(flight.Destination))
		}

		// Process source state (outgoing flights)
		if sourceOk {
			sourceState = strings.Title(strings.ToLower(sourceState)) // Capitalize properly
			if _, exists := aggregations[sourceState]; !exists {
				aggregations[sourceState] = &StateAggregation{
					StateName:       sourceState,
					TotalFlights:    0,
					IncomingFlights: 0,
					OutgoingFlights: 0,
					UniqueRoutes:    0,
					Airlines:        make(map[string]int),
					RouteDetails:    make(map[string]int),
				}
			}

			agg := aggregations[sourceState]
			agg.OutgoingFlights++
			agg.TotalFlights++
			agg.Airlines[flight.Airline]++

			// Add route detail
			routeKey := strings.ToLower(flight.Source + "->" + flight.Destination)
			agg.RouteDetails[routeKey]++
		}

		// Process destination state (incoming flights)
		if destOk {
			destState = strings.Title(strings.ToLower(destState)) // Capitalize properly
			if _, exists := aggregations[destState]; !exists {
				aggregations[destState] = &StateAggregation{
					StateName:       destState,
					TotalFlights:    0,
					IncomingFlights: 0,
					OutgoingFlights: 0,
					UniqueRoutes:    0,
					Airlines:        make(map[string]int),
					RouteDetails:    make(map[string]int),
				}
			}

			agg := aggregations[destState]
			agg.IncomingFlights++
			agg.TotalFlights++
			agg.Airlines[flight.Airline]++

			// Add route detail (same route key for both source and dest aggregations)
			routeKey := strings.ToLower(flight.Source + "->" + flight.Destination)
			agg.RouteDetails[routeKey]++
		}

	}

	// Calculate unique routes for each state
	for _, agg := range aggregations {
		agg.UniqueRoutes = len(agg.RouteDetails)
	}

	sa.aggregations = aggregations

	log.Printf("Computed state-wise aggregations for %d states", len(aggregations))

	// Log some summary information
	for state, agg := range aggregations {
		log.Printf("State: %s - Total: %d, Incoming: %d, Outgoing: %d, Unique Routes: %d, Airlines: %d",
			state, agg.TotalFlights, agg.IncomingFlights, agg.OutgoingFlights, agg.UniqueRoutes, len(agg.Airlines))
	}
}

// Returns the aggregation and a boolean indicating if it exists
func (sa *StateAggregator) GetAggregationForState(stateName string) (*StateAggregation, bool) {
	sa.mutex.RLock()
	defer sa.mutex.RUnlock()

	// Normalize state name for lookup (handle different cases)
	normalizedState := strings.Title(strings.ToLower(stateName))

	agg, exists := sa.aggregations[normalizedState]
	if exists {
		return agg, true
	}

	// Check if this is a valid Indian state name, even if it has no flight data
	allStates := sa.GetAllIndianStates()
	for _, validState := range allStates {
		if strings.EqualFold(validState, stateName) || strings.EqualFold(strings.Title(strings.ToLower(validState)), normalizedState) {
			// Return a default aggregation with 0 values for valid states without data
			defaultAgg := &StateAggregation{
				StateName:       validState,
				TotalFlights:    0,
				IncomingFlights: 0,
				OutgoingFlights: 0,
				UniqueRoutes:    0,
				Airlines:        make(map[string]int),
				RouteDetails:    make(map[string]int),
			}
			return defaultAgg, true
		}
	}

	return nil, false
}

// GetAllAggregations returns all computed state aggregations
func (sa *StateAggregator) GetAllAggregations() map[string]*StateAggregation {
	sa.mutex.RLock()
	defer sa.mutex.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[string]*StateAggregation)
	for k, v := range sa.aggregations {
		result[k] = v
	}
	return result
}

// RefreshAggregations recomputes all aggregations (useful when data changes)
func (sa *StateAggregator) RefreshAggregations() {
	log.Println("Refreshing state-wise aggregations...")
	sa.ComputeAggregations()
}

// GetAllIndianStates returns a list of all Indian states
func (sa *StateAggregator) GetAllIndianStates() []string {
	allStates := []string{
		"Andhra Pradesh", "Arunachal Pradesh", "Assam", "Bihar", "Chhattisgarh",
		"Goa", "Gujarat", "Haryana", "Himachal Pradesh", "Jharkhand",
		"Karnataka", "Kerala", "Madhya Pradesh", "Maharashtra", "Manipur",
		"Meghalaya", "Mizoram", "Nagaland", "Odisha", "Punjab",
		"Rajasthan", "Sikkim", "Tamil Nadu", "Telangana", "Tripura",
		"Uttar Pradesh", "Uttarakhand", "West Bengal",
		"Delhi", "Puducherry", "Andaman and Nicobar Islands",
		"Dadra and Nagar Haveli and Daman and Diu", "Lakshadweep", "Ladakh",
	}
	return allStates
}

// GetStatesList returns a list of all states that have flight data
func (sa *StateAggregator) GetStatesList() []string {
	sa.mutex.RLock()
	defer sa.mutex.RUnlock()

	states := make([]string, 0, len(sa.aggregations))
	for state := range sa.aggregations {
		states = append(states, state)
	}
	return states
}

// GetTopAirlinesForState returns the top airlines for a specific state
func (sa *StateAggregator) GetTopAirlinesForState(stateName string, limit int) map[string]int {
	agg, exists := sa.GetAggregationForState(stateName)
	if !exists {
		return nil
	}

	// Return a copy of the airlines map
	result := make(map[string]int)
	for k, v := range agg.Airlines {
		result[k] = v
	}

	// In a real implementation, we would sort and limit the results
	// For now, returning the full map
	return result
}

// GetTotalFlightsForState returns the total number of flights for a specific state
func (sa *StateAggregator) GetTotalFlightsForState(stateName string) int {
	agg, exists := sa.GetAggregationForState(stateName)
	if !exists {
		return 0
	}
	return agg.TotalFlights
}

// GetIncomingFlightsForState returns the number of incoming flights for a specific state
func (sa *StateAggregator) GetIncomingFlightsForState(stateName string) int {
	agg, exists := sa.GetAggregationForState(stateName)
	if !exists {
		return 0
	}
	return agg.IncomingFlights
}

// GetOutgoingFlightsForState returns the number of outgoing flights for a specific state
func (sa *StateAggregator) GetOutgoingFlightsForState(stateName string) int {
	agg, exists := sa.GetAggregationForState(stateName)
	if !exists {
		return 0
	}
	return agg.OutgoingFlights
}

// GetUniqueRoutesForState returns the number of unique routes for a specific state
func (sa *StateAggregator) GetUniqueRoutesForState(stateName string) int {
	agg, exists := sa.GetAggregationForState(stateName)
	if !exists {
		return 0
	}
	return agg.UniqueRoutes
}

// normalizeCityNameForMapping normalizes city names to match the mapper's expected format
func normalizeCityNameForMapping(city string) string {
	return strings.ToLower(strings.TrimSpace(city))
}

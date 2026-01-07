package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"flight-dashboard-backend/services"

	"github.com/labstack/echo/v4"
)

// GetStateWiseFlights returns state-wise flight aggregations
func GetStateWiseFlights(c echo.Context) error {
	aggregator := services.GetStateAggregator()

	// Get state from query parameter if provided
	stateParam := c.QueryParam("state")

	if stateParam != "" {
		// Return aggregation for specific state
		agg, exists := aggregator.GetAggregationForState(stateParam)
		if !exists {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "State not found: " + stateParam,
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"data":    agg,
		})
	} else {
		// Return all state aggregations
		allAggs := aggregator.GetAllAggregations()
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"data":    allAggs,
			"count":   len(allAggs),
		})
	}
}

// GetStateList returns a list of all states with flight data summary
func GetStateList(c echo.Context) error {
	aggregator := services.GetStateAggregator()
	allAggs := aggregator.GetAllAggregations()

	// Get all Indian states
	allIndianStates := aggregator.GetAllIndianStates()

	// Format response as array of state objects with total flights
	stateSummaries := make([]map[string]interface{}, 0, len(allIndianStates))
	for _, stateName := range allIndianStates {
		// Check if this state has flight data
		normalizedStateName := strings.Title(strings.ToLower(stateName))
		if agg, exists := allAggs[normalizedStateName]; exists {
			stateSummary := map[string]interface{}{
				"state":        stateName,
				"totalFlights": agg.TotalFlights,
			}
			stateSummaries = append(stateSummaries, stateSummary)
		} else {
			// Include states with 0 flights
			stateSummary := map[string]interface{}{
				"state":        stateName,
				"totalFlights": 0,
			}
			stateSummaries = append(stateSummaries, stateSummary)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    stateSummaries,
		"count":   len(stateSummaries),
	})
}

// GetStateDetail returns detailed information for a specific state
// Provides state-wise flight statistics for hover functionality
// Response format: {"state": "Karnataka", "totalFlights": 2100, "incomingFlights": 980, "outgoingFlights": 1120, "routes": 120, "airlines": ["IndiGo", "Vistara", "Air India"]}
func GetStateDetail(c echo.Context) error {
	stateParam := c.Param("state")

	// Convert kebab-case to proper case (e.g., "rajasthan" -> "Rajasthan")
	normalizedState := normalizeStateName(stateParam)

	aggregator := services.GetStateAggregator()
	agg, exists := aggregator.GetAggregationForState(normalizedState)

	// If not found with normalized name, try the original parameter
	if !exists {
		agg, exists = aggregator.GetAggregationForState(stateParam)
	}

	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "State not found: " + stateParam,
		})
	}

	// Extract airline names from the map
	airlines := make([]string, 0, len(agg.Airlines))
	for airline := range agg.Airlines {
		airlines = append(airlines, airline)
	}

	// Create response object matching the required format
	response := map[string]interface{}{
		"state":           agg.StateName,
		"totalFlights":    agg.TotalFlights,
		"incomingFlights": agg.IncomingFlights,
		"outgoingFlights": agg.OutgoingFlights,
		"routes":          agg.UniqueRoutes,
		"airlines":        airlines,
	}

	return c.JSON(http.StatusOK, response)
}

// normalizeStateName converts kebab-case state names to proper title case
func normalizeStateName(state string) string {
	// Convert kebab-case to space-separated
	state = strings.ReplaceAll(state, "-", " ")
	// Convert to title case (first letter of each word capitalized)
	state = strings.Title(state)
	// Handle special cases like "And" and "Or" that should be lowercase in state names
	words := strings.Split(state, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if lowerWord == "and" || lowerWord == "or" {
			words[i] = lowerWord
		}
	}
	return strings.Join(words, " ")
}

// GetTopAirlinesForState returns top airlines for a specific state
func GetTopAirlinesForState(c echo.Context) error {
	state := c.Param("state")
	limitStr := c.QueryParam("limit")

	limit := 10 // default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	aggregator := services.GetStateAggregator()
	airlines := aggregator.GetTopAirlinesForState(state, limit)

	if airlines == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "State not found: " + state,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"state":   state,
		"data":    airlines,
		"count":   len(airlines),
	})
}

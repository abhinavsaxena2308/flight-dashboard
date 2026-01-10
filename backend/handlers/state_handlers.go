package handlers

import (
	"flight-dashboard-backend/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// gets flight data grouped by state - pretty useful for the dashboard
func GetStateWiseFlights(c echo.Context) error {
	aggregator := services.GetStateAggregator()

	// get state from query parameter if provided
	stateParam := c.QueryParam("state")

	if stateParam != "" {
		// gives aggregation for specific state
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
		// gives all state aggregations
		allAggs := aggregator.GetAllAggregations()
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"data":    allAggs,
			"count":   len(allAggs),
		})
	}
}

// returns a list of all states with basic flight data - needed for the state selection page
func GetStateList(c echo.Context) error {
	aggregator := services.GetStateAggregator()
	allAggs := aggregator.GetAllAggregations()

	// getting all states
	allIndianStates := aggregator.GetAllIndianStates()

	// formatting response as array of state objects with total flights
	stateSummaries := make([]map[string]interface{}, 0, len(allIndianStates))
	for _, stateName := range allIndianStates {
		// checking if a state is having flight data or not
		normalizedStateName := strings.Title(strings.ToLower(stateName))
		if agg, exists := allAggs[normalizedStateName]; exists {
			stateSummary := map[string]interface{}{
				"state":        stateName,
				"totalFlights": agg.TotalFlights,
			}
			stateSummaries = append(stateSummaries, stateSummary)
		} else {
			// including states with 0 flights
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


// getting detailed info of state ..... this is used for the hover tooltips on the map
// returns data in the format: {"state": "Karnataka", "totalFlights": 2100, "incomingFlights": 980, "outgoingFlights": 1120, "routes": 120, "airlines": ["IndiGo", "Vistara", "Air India"]}
func GetStateDetail(c echo.Context) error {
	stateParam := c.Param("state")
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

	// retriving airline names from the map
	airlines := make([]string, 0, len(agg.Airlines))
	for airline := range agg.Airlines {
		airlines = append(airlines, airline)
	}

	// response format
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

// converts kebab-case state names (like 'tamil-nadu') to proper format (like 'Tamil Nadu')
func normalizeStateName(state string) string {
	state = strings.ReplaceAll(state, "-", " ")
	state = strings.Title(state)
	words := strings.Split(state, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if lowerWord == "and" || lowerWord == "or" {
			words[i] = lowerWord
		}
	}
	return strings.Join(words, " ")
}

// returns the top airlines for a specific state - used for the airline breakdown section
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

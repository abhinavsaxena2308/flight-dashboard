# Flight Dashboard Backend

A Go-based backend service for the Indian flight data dashboard with state-wise aggregation and visualization capabilities.

## Overview

This backend service provides APIs to access and analyze Indian domestic flight data, with a focus on state-wise aggregation for visualization on an interactive India map. The service processes flight data from the Goibibo dataset and provides endpoints for various flight statistics.

## Features

- **State-wise Flight Aggregation**: Precomputed aggregations for all Indian states
- **Fast API Responses**: Optimized for real-time visualization (especially hover interactions)
- **City-to-State Mapping**: Comprehensive mapping of Indian cities to their respective states
- **CSV Data Loading**: Automatic loading and processing of flight data at startup
- **RESTful API Endpoints**: Clean and well-documented API for frontend integration

## API Endpoints

### State Summary
- `GET /api/states` - Get summary of all states with flight data
  - Response: `[{ "state": "Maharashtra", "totalFlights": 3450 }]`

### State Detail (Hover API)
- `GET /api/state/{stateName}` - Get detailed information for a specific state
  - Response: 
    ```json
    {
      "state": "Karnataka",
      "totalFlights": 2100,
      "incomingFlights": 980,
      "outgoingFlights": 1120,
      "routes": 120,
      "airlines": ["IndiGo", "Vistara", "Air India"]
    }
    ```

### Additional Endpoints
- `GET /api/state-flights` - Get all state-wise flight aggregations
- `GET /api/states/:state/airlines` - Get airlines for a specific state
- `GET /health` - Health check endpoint

## Architecture

### Data Flow
1. **Data Loading**: CSV flight data is loaded at application startup
2. **State Mapping**: City names are mapped to Indian states using comprehensive city-to-state mapping
3. **Aggregation**: All state-wise metrics are precomputed for fast access
4. **API Access**: Endpoints provide access to aggregated data with minimal processing

### Performance Optimizations
- **Precomputed Aggregations**: All state-wise statistics are calculated once at startup
- **In-Memory Storage**: Aggregated data is stored in memory for O(1) access
- **Thread-Safe Operations**: RWMutex used for concurrent access to data
- **Efficient Data Structures**: Maps used for fast lookups

## Project Structure

```
backend/
├── main.go                 # Application entry point
├── routes/
│   └── routes.go          # API route definitions
├── handlers/
│   ├── health_handler.go  # Health check handler
│   └── state_handlers.go  # State-wise data handlers
├── services/
│   ├── flight_data_service.go    # Flight data loading and access
│   ├── city_state_mapper.go      # City-to-state mapping service
│   └── state_aggregator.go       # State-wise aggregation service
└── data/
    ├── dataset.csv        # Flight data CSV file
    └── city_state_map.json # City-to-state mapping data
```

## Setup

### Prerequisites
- Go 1.19 or higher
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd flight-dashboard/backend
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Place your flight data CSV file as `data/dataset.csv`

4. Run the application:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

## Configuration

The application is configured via the source code. Key configuration points:
- **Port**: Set in `main.go` (default: 8080)
- **Data file**: Set in `main.go` (default: `data/dataset.csv`)
- **CORS**: Enabled by default for frontend integration

## Performance Notes

- The application precomputes all state-wise aggregations at startup for fast API responses
- State detail endpoints (used for hover functionality) respond in milliseconds
- Memory usage scales linearly with the size of the flight dataset
- The city-to-state mapping handles common variations in city names (e.g., Bombay → Mumbai)

## Data Processing

The service handles common data quality issues:
- Normalizes city names (e.g., "bombay" → "mumbai")
- Handles various duration formats (e.g., "2h 15m")
- Processes price values with commas (e.g., "5,000")
- Manages different stop formats (e.g., "non-stop", "1-stop")

## Error Handling

- Graceful handling of missing or invalid city names
- Proper error responses for non-existent states
- Logging for debugging and monitoring
- Safe fallbacks when data is missing

## Usage with Frontend

The API is designed to work seamlessly with the frontend India map visualization:
- The `/api/states` endpoint provides initial data for map display
- The `/api/state/{stateName}` endpoint provides hover details with fast response times
- All responses are in JSON format suitable for JavaScript processing

## Development

To run tests:
```bash
go test ./...
```

To build for production:
```bash
go build -o flight-dashboard main.go
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

[Specify your license here]
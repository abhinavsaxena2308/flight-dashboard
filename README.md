# Indian Flight Dashboard

An interactive dashboard for visualizing Indian domestic flight data on an interactive India map with state-wise aggregations.

## Overview

This project consists of:
- **Frontend**: A Next.js application with an interactive India map visualization
- **Backend**: A Go-based REST API that processes flight data and provides state-wise aggregations

## Features

- Interactive India map with state highlighting
- Real-time flight statistics on hover
- State-wise flight aggregations (total, incoming, outgoing flights)
- Route and airline information per state
- Responsive design for different screen sizes

## Tech Stack

### Frontend
- Next.js 14
- React
- TypeScript
- Leaflet (for map visualization)
- react-leaflet (React bindings for Leaflet)
- Tailwind CSS

### Backend
- Go (Golang)
- Echo framework
- encoding/csv (for data processing)

## Prerequisites

- Node.js (v16 or higher)
- Go (v1.19 or higher)
- npm or yarn

## Installation & Setup

### Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

3. Place your flight data CSV file as `data/dataset.csv` (or use the sample data provided)

4. Run the backend server:
   ```bash
   go run main.go
   ```
   The backend will start on `http://localhost:8080`

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Run the development server:
   ```bash
   npm run dev
   ```
   The frontend will start on `http://localhost:3000`

## Running Both Servers Together

To run both frontend and backend servers simultaneously:

```bash
npm run dev:fullstack
```

This command will start both the backend server on port 8080 and the frontend server on port 3000, with API requests automatically proxied from the frontend to the backend.

## API Endpoints

The backend provides the following API endpoints:

- `GET /api/states` - Get summary of all states with flight data
  - Response: `[{ "state": "Maharashtra", "totalFlights": 3450 }]`

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

- `GET /health` - Health check endpoint

## How It Works

1. The backend loads flight data from CSV at startup and precomputes state-wise aggregations
2. When a user hovers over a state in the frontend map, an API call is made to `/api/state/{stateName}`
3. The backend responds with detailed statistics for that state
4. The frontend displays this information in a tooltip near the cursor

## Project Structure

```
project-root/
├── backend/                 # Go backend server
│   ├── main.go             # Entry point
│   ├── routes/             # API route definitions
│   ├── handlers/           # Request handlers
│   ├── services/           # Business logic
│   │   ├── flight_data_service.go
│   │   ├── city_state_mapper.go
│   │   └── state_aggregator.go
│   └── data/               # Data files
│       ├── dataset.csv     # Flight data
│       └── city_state_map.json
└── frontend/               # Next.js frontend
    ├── src/
    │   ├── app/            # Next.js app directory
    │   ├── components/     # React components
    │   │   ├── IndiaMap.tsx
    │   │   └── MapWithStates.tsx
    │   └── styles/
    ├── public/
    ├── next.config.js      # Next.js configuration with API proxy
    └── package.json
```

## Configuration

### Frontend Configuration

The `next.config.js` file includes a rewrite rule that proxies API requests:
```js
async rewrites() {
  return [
    {
      source: '/api/:path*',
      destination: 'http://localhost:8080/api/:path*',
    },
  ];
},
```

This allows the frontend to make API calls to `/api/` which get forwarded to the backend server.

### Backend Configuration

- Port: 8080 (configured in `main.go`)
- CORS enabled by default
- Data file: `data/dataset.csv`

## Development

To run the full stack during development:

1. Make sure both backend (port 8080) and frontend (port 3000) are running
2. The frontend will automatically proxy API requests to the backend
3. Changes to frontend code will hot-reload
4. For backend changes, restart the Go server

## Building for Production

### Backend
```bash
cd backend
go build -o flight-dashboard main.go
./flight-dashboard
```

### Frontend
```bash
cd frontend
npm run build
npm start
```

## Troubleshooting

### Common Issues

1. **API calls failing**: Ensure the backend server is running on port 8080
2. **Map not loading**: Check browser console for Leaflet-related errors
3. **Data not displaying**: Verify that `data/dataset.csv` exists in the backend directory

### Debugging API Calls

API requests from the frontend are logged by the backend. Check the backend console for:
- Request paths
- Response times
- Error messages

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

[Specify your license here]
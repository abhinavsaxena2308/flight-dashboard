package models
type Flight struct {
	Airline        string  `json:"airline"`
	FlightDate     string  `json:"flight_date"`
	Source         string  `json:"source"`
	Destination    string  `json:"destination"`
	FlightClass    string  `json:"flight_class"`
	Duration       float64 `json:"duration"`
	Price          float64 `json:"price"`
	DepartureTime  string  `json:"departure_time"`
	ArrivalTime    string  `json:"arrival_time"`
	Stops          int     `json:"stops"`
	AdditionalInfo string  `json:"additional_info"`
}

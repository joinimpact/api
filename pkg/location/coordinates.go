package location

// Coordinates represents a latitude and longitude representation of a location.
type Coordinates struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
}

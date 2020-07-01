package location

// Coordinates represents a latitude and longitude representation of a location.
type Coordinates struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"long"`
}

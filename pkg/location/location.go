package location

// Location represents the name of a single location, broken up into
// different subsections.
type Location struct {
	City    *LocationName `json:"city"`
	State   *LocationName `json:"state"`
	Country *LocationName `json:"country"`
}

// LocationName contains a short and long name for a location.
type LocationName struct {
	ShortName string `json:"shortName,omitempty"`
	LongName  string `json:"longName,omitempty"`
}

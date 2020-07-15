package location

// Location represents the name of a single location, broken up into
// different subsections.
type Location struct {
	StreetAddress *LocationName `json:"streetAddress,omitempty"`
	City          *LocationName `json:"city,omitempty"`
	State         *LocationName `json:"state,omitempty"`
	Country       *LocationName `json:"country,omitempty"`
}

// LocationName contains a short and long name for a location.
type LocationName struct {
	ShortName string `json:"shortName,omitempty"`
	LongName  string `json:"longName,omitempty"`
}

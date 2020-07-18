package opportunities

// OpportunityDocument represents an opportunity as an Elasticsearch NoSQL document.
type OpportunityDocument struct {
	ID           int64                            `json:"opportunityId"`
	Public       bool                             `json:"public"`
	Organization *OpportunityOrganizationDocument `json:"organization"`
	Title        string                           `json:"title"`
	Description  string                           `json:"description"`
	Tags         []OpportunityTagDocument         `json:"tags"`
	Location     *LocationDocument                `json:"location,omitempty"`
	Requirements *Requirements                    `json:"requirements"`
	Limits       *Limits                          `json:"limits"`
}

// LocationDocument contains a location.
type LocationDocument struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

// OpportunityOrganizationDocument has a summary of the organization for an
// opportunity.
type OpportunityOrganizationDocument struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// OpportunityTagDocument represents an opportunity's tag in the Elasticsearch database.
type OpportunityTagDocument struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Category int    `json:"category"`
}

// Requirements defines the requirements to join an opportunity.
type Requirements struct {
	AgeLimit      AgeLimit      `json:"ageLimit"`
	ExpectedHours ExpectedHours `json:"expectedHours"`
}

// AgeLimit represents age limit requirements.
type AgeLimit struct {
	Active bool `json:"active"`
	From   int  `json:"from"`
	To     int  `json:"to"`
}

// ExpectedHours represents expected weekly hour requirements.
type ExpectedHours struct {
	Active bool `json:"active"`
	Hours  int  `json:"hours"`
}

// Limits defines the limits of an opportunity.
type Limits struct {
	VolunteersCap VolunteersCap `json:"volunteersCap"`
}

// VolunteersCap limits how many volunteers can join.
type VolunteersCap struct {
	Active bool `json:"active"`
	Cap    int  `json:"cap"`
}

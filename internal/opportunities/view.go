package opportunities

// OpportunityView is a representation of how the opportunities will be returned by the API.
type OpportunityView struct {
	ID             int64         `json:"id"`
	OrganizationID int64         `json:"organizationId"`
	CreatorID      int64         `json:"creatorId" scope:"manager"`
	ProfilePicture string        `json:"profilePicture"`
	Title          string        `json:"title"`
	Description    string        `json:"description"`
	Public         bool          `json:"public" scope:"manager"`
	Stats          *Stats        `json:"stats" scope:"manager"`
	Requirements   *Requirements `json:"requirements"`
	Limits         *Limits       `json:"limits"`
}

// Stats represents opportunity statistics.
type Stats struct {
	VolunteersEnrolled int `json:"volunteersEnrolled"`
	VolunteersPending  int `json:"volunteersPending"`
}

// Requirements defines the requirements to join an opportunity.
type Requirements struct {
	AgeLimit      AgeLimit      `json:"ageLimit"`
	ExpectedHours ExpectedHours `json:"expectedHours"`
}

// AgeLimit represents age limit requirements.
type AgeLimit struct {
	Active bool `json:"active"`
	From   int  `json:"from,omitempty"`
	To     int  `json:"to,omitempty"`
}

// ExpectedHours represents expected weekly hour requirements.
type ExpectedHours struct {
	Active bool `json:"active"`
	Hours  int  `json:"hours,omitempty"`
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

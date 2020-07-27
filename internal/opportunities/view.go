package opportunities

import (
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/pkg/location"
)

// OpportunityView is a representation of how the opportunities will be returned by the API.
type OpportunityView struct {
	ID                             int64                           `json:"id"`
	Publishable                    bool                            `json:"publishable"`
	OrganizationID                 int64                           `json:"organizationId"`
	CreatorID                      int64                           `json:"creatorId" scope:"manager"`
	ProfilePicture                 string                          `json:"profilePicture"`
	Title                          string                          `json:"title"`
	Description                    string                          `json:"description"`
	Location                       *location.Location              `json:"location" validate:"-"`
	Public                         bool                            `json:"public" scope:"authenticated"`
	Tags                           []models.Tag                    `json:"tags"` // the model's tags
	Stats                          *Stats                          `json:"stats" scope:"authenticated"`
	Requirements                   *Requirements                   `json:"requirements"`
	Limits                         *Limits                         `json:"limits"`
	OpportunityOrganizationProfile *OpportunityOrganizationProfile `json:"organization,omitempty"`
}

// OpportunityOrganizationProfile contains an organization profile in an opportunity.
type OpportunityOrganizationProfile struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profilePicture"`
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

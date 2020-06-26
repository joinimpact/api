package users

import (
	"time"

	"github.com/joinimpact/api/internal/models"
)

// UserProfile represents a user's profile.
type UserProfile struct {
	ProfilePicture string       `json:"profilePicture,omitempty"`        // a URL for the user's profile picture
	FirstName      string       `json:"firstName"`                       // the user's first name
	LastName       string       `json:"lastName"`                        // the user's last name
	Email          string       `json:"email,omitempty"`                 // the user's email
	DateOfBirth    time.Time    `json:"dateOfBirth,omitempty" level:"1"` // the user's date of birth, used for calculating age
	ZIPCode        string       `json:"zipCode,omitempty" level:"1"`     // the user's zip code, used to find nearby opportunities
	Tags           []models.Tag `json:"tags"`
}

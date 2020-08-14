package users

import (
	"time"

	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/pkg/location"
)

// UserProfile represents a user's profile.
type UserProfile struct {
	ID             int64                     `json:"id"`
	ProfilePicture string                    `json:"profilePicture,omitempty"`            // a URL for the user's profile picture
	FirstName      string                    `json:"firstName"`                           // the user's first name
	LastName       string                    `json:"lastName"`                            // the user's last name
	Email          string                    `json:"email,omitempty" scope:"owner"`       // the user's email
	LastOnline     *time.Time                `json:"lastOnline,omitempty"`                // the last time the user was online
	DateOfBirth    time.Time                 `json:"dateOfBirth,omitempty" scope:"owner"` // the user's date of birth, used for calculating age
	CreatedAt      time.Time                 `json:"createdAt,omitempty" scope:"owner"`   // the time the user was created at
	Tags           []models.Tag              `json:"tags"`                                // the user's tags
	Location       *location.Location        `json:"location,omitempty" scope:"owner"`    // a formatted location
	ProfileFields  []models.UserProfileField `json:"profile"`                             // the user's profile fields
}

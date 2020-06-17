package models

import (
	"time"
)

// User represents a single user in the Impact application.
type User struct {
	Model
	Active         bool               `json:"-"`                     // controls whether or not the account is active (false if the account is suspended)
	Email          string             `json:"email" level:"1"`       // the user's email address
	EmailVerified  bool               `json:"emailVerified"`         // whether or not the user has verified their email
	Password       string             `json:"-"`                     // a bcrypt hash of the user's passowrd
	ProfilePicture string             `json:"profilePicture"`        // a URL for the user's profile picture
	FirstName      string             `json:"firstName"`             // the user's first name
	LastName       string             `json:"lastName"`              // the user's last name
	DateOfBirth    time.Time          `json:"dateOfBirth" level:"1"` // the user's date of birth, used for calculating age
	ZIPCode        string             `json:"zipCode" level:"1"`     // the user's zip code, used to find nearby opportunities
	ProfileFields  []UserProfileField `json:"profile"`               // fields of the user's profile
}

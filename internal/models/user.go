package models

import (
	"time"
)

// User represents a single user in the Impact application.
type User struct {
	Model
	Active         bool               `json:"-" gorm:"default:'true'"`     // controls whether or not the account is active (false if the account is suspended)
	Email          string             `json:"email" scope:"owner"`         // the user's email address
	EmailVerified  bool               `json:"emailVerified" scope:"owner"` // whether or not the user has verified their email
	Password       string             `json:"-"`                           // a bcrypt hash of the user's passowrd
	ProfilePicture string             `json:"profilePicture"`              // a URL for the user's profile picture
	FirstName      string             `json:"firstName"`                   // the user's first name
	LastName       string             `json:"lastName"`                    // the user's last name
	DateOfBirth    time.Time          `json:"dateOfBirth" level:"1"`       // the user's date of birth, used for calculating age
	LastOnline     *time.Time         `json:"lastOnline"`                  // the last time the user was online
	ProfileFields  []UserProfileField `json:"profile"`                     // fields of the user's profile
	// ZIPCode        string             `json:"zipCode" level:"1"`       // the user's zip code, used to find nearby opportunities
	LocationLatitude  float64 `json:"-"` // the latitude of the user's city
	LocationLongitude float64 `json:"-"` // the longitude of the user's city
}

// UserRepository represents a repository of users.
type UserRepository interface {
	// FindByID finds a single User by ID.
	FindByID(id int64) (*User, error)
	// FindByEmail finds a single User by Email.
	FindByEmail(email string) (*User, error)
	// Create creates a new User.
	Create(user User) error
	// Update updates a User with the ID in the provided User.
	Update(user User) error
	// DeleteByID deletes a User by ID.
	DeleteByID(id int64) error
}

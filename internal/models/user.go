package models

import "github.com/jinzhu/gorm"

// User represents a single user in the Impact application.
type User struct {
	gorm.Model
	ID            int64          `gorm:"primary_key" json:"id"` // the user's unique int64 ID
	Active        bool           `json:"-"`                     // controls whether or not the account is active (false if the account is suspended)
	Email         string         `json:"email"`                 // the user's email address
	EmailVerified bool           `json:"verified"`              // whether or not the user has verified their email
	FirstName     string         `json:"firstName"`             // the user's first name
	LastName      string         `json:"lastName"`              // the user's last name
	Password      string         `json:"-"`                     // a bcrypt hash of the user's passowrd
	ProfileFields []ProfileField `json:"profile"`               // fields of the user's profile
}

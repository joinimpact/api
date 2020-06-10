package models

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

// UserProfileField represents a single field in a user's profile.
type UserProfileField struct {
	gorm.Model
	ID     int64  `gorm:"primary_key" json:"id"`
	UserID int64  `json:"userId"`
	Name   string `json:"name"`            // the key/name of the profile field
	Value  string `json:"value,omitempty"` // the value of the profile field (if it is a string)
	// School School
	ValueInt  int            `json:"value,omitempty"` // the value of the profile field (if it is an int)
	ValueJSON postgres.Jsonb `json:"value,omitempty"` // the value of the profile field in JSON
	MixedID   int64          `json:"-"`               // represents an ID for a relation to a mixed-type item in the database.
}

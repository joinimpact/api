package models

import "github.com/jinzhu/gorm"

// ProfileField represents a single field in a user's profile.
type ProfileField struct {
	gorm.Model
	ID    int64  `gorm:"primary_key" json:"id"`
	Name  string `json:"name"`            // the key/name of the profile field
	Value string `json:"value,omitempty"` // the value of the profile field (if it is a string)
	// School School
	ValueInt int   `json:"value,omitempty"` // the value of the profile field (if it is an int)
	MixedID  int64 `json:"-"`               // represents an ID for a relation to a mixed-type item in the database.
}

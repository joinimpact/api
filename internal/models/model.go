package models

import "time"

// Model is the base unit for each database model.
type Model struct {
	ID        int64 `gorm:"primary_key" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

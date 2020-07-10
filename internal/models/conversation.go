package models

// Conversation represents a single chat conversation in the database.
type Conversation struct {
	Model
	Active         bool  `json:"-"`
	CreatorID      int64 `json:"creatorId"`
	OrganizationID int64 `json:"organizationID"`
	Organization   Organization
	Type           int `json:"type"`
}

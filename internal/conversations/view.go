package conversations

import (
	"time"

	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/pkg/location"
)

// ConversationView represents a view of a conversation.
type ConversationView struct {
	models.Conversation
	LastMessageView *MessageView `json:"lastMessage"`
	UnreadCount     int          `json:"unreadCount"`
}

// MessageView represents a view of a message.
type MessageView struct {
	ID                int64              `json:"id"`
	ConversationID    int64              `json:"conversationId"`
	SenderID          int64              `json:"senderId"`
	SenderPerspective uint               `json:"senderPerspective"`
	Timestamp         time.Time          `json:"timestamp"`
	Type              string             `json:"type"`
	Edited            bool               `json:"edited"`
	EditedTimestamp   time.Time          `json:"editedTimestamp"`
	Body              interface{}        `json:"body"`
	Sender            *MessageSenderView `json:"sender"`
}

// MessageSenderView represents the sender of a message.
type MessageSenderView struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	ProfilePicture string `json:"profilePicture"`
}

// MessageVolunteerRequestProfileView represents a view of a message containing a user's profile.
type MessageVolunteerRequestProfileView struct {
	UserID             int64                     `json:"userId"`
	ProfilePicture     string                    `json:"profilePicture,omitempty"`            // a URL for the user's profile picture
	FirstName          string                    `json:"firstName"`                           // the user's first name
	LastName           string                    `json:"lastName"`                            // the user's last name
	DateOfBirth        time.Time                 `json:"dateOfBirth,omitempty" scope:"owner"` // the user's date of birth, used for calculating age
	PreviousExperience *PreviousExperience       `json:"previousExperience"`
	Tags               []models.Tag              `json:"tags"`                             // the user's tags
	Location           *location.Location        `json:"location,omitempty" scope:"owner"` // a formatted location
	ProfileFields      []models.UserProfileField `json:"profile"`                          // the user's profile fields
	Message            string                    `json:"message"`
}

// PreviousExperience represents a user's previous experience.
type PreviousExperience struct {
	Count int `json:"count"`
}

// MessageTypeVolunteerRequestAcceptanceView represents a view of a message containing an opportunity acceptance.
type MessageTypeVolunteerRequestAcceptanceView struct {
	Volunteer        *MessageUserWithName `json:"volunteer"`
	Accepter         *MessageUserWithName `json:"accepter"`
	OpportunityID    int64                `json:"opportunityId"`
	OpportunityTitle string               `json:"opportunityTitle"`
}

// MessageUserWithName represents a first and last name pair.
type MessageUserWithName struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// MessageTypeHoursRequestedView represents the message sent when a volunteer requests hours from an organization.
type MessageTypeHoursRequestedView struct {
	models.VolunteeringHourLogRequest
}

// MessageTypeHoursAcceptedView represents the message sent when a volunteer's request is accepted.
type MessageTypeHoursAcceptedView struct {
	models.VolunteeringHourLogRequest
}

// MessageTypeHoursDeclinedView represents the message sent when a volunteer's request is declined.
type MessageTypeHoursDeclinedView struct {
	models.VolunteeringHourLogRequest
}

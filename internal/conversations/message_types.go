package conversations

// MessageStandard represents a standard message.
type MessageStandard struct {
	Text string `json:"text"`
}

// MessageVolunteerRequestProfile represents the message sent containing a user's profile when they request to join an opportunity.
type MessageVolunteerRequestProfile struct {
	Message string `json:"message"`
	UserID  int64  `json:"userId"`
}

// MessageTypeVolunteerRequestAcceptance represents the message sent when a user is accepted to an opportunity.
type MessageTypeVolunteerRequestAcceptance struct {
	UserID        int64 `json:"userId"`
	OpportunityID int64 `json:"opportunityId"`
}

// MessageTypeHoursRequested represents the message sent when a volunteer requests hours from an organization.
type MessageTypeHoursRequested struct {
	VolunteeringHourLogRequestID int64 `json:"requestId"`
}

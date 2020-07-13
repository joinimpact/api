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

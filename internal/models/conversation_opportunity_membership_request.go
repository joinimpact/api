package models

// ConversationOpportunityMembershipRequest represents a connection of a conversation to a single opportunity membership request.
type ConversationOpportunityMembershipRequest struct {
	Model
	ConversationID                 int64
	OpportunityMembershipRequestID int64
}

package models

import "context"

const (
	OpportunityApplicationQuestionTypeShortAnswer uint = iota
	OpportunityApplicationQuestionTypeParagraph   uint = iota
)

// OppOpportunityApplicationQuestion represents a question to be answered by a volunteer when applying for an opportunity.
type OpportunityApplicationQuestion struct {
	Model
	Opportunity   Opportunity `json:"-"`
	OpportunityID int64       `json:"opportunityId"`
	Required      *bool       `json:"required"`
	Type          *uint       `json:"type"`
	Title         string      `json:"title"`
	MinLength     *uint       `json:"minLength"`
	MaxLength     *uint       `json:"maxLength"`
}

// OpportunityApplicationQuestionsResponse represents a response containing multiple OpportunityApplicationQuestions.
type OpportunityApplicationQuestionsResponse struct {
	OpportunityApplicationQuestions []OpportunityApplicationQuestion
	TotalResults                    int
}

// OpportunityOpportunityApplicationQuestionRepository defines an interface of methods for interacting with OpportunityApplicationQuestion entities.
type OpportunityApplicationQuestionRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(ctx context.Context, id int64) (*OpportunityApplicationQuestion, error)
	// FindByOpportunityID finds multiple entities by OpportunityID.
	FindByOpportunityID(ctx context.Context, opportunityID int64) (*OpportunityApplicationQuestionsResponse, error)
	// Create creates a new entity.
	Create(ctx context.Context, opportunityApplicationQuestion OpportunityApplicationQuestion) error
	// Update updates an entity with the ID in the provided entity.
	Update(ctx context.Context, opportunityApplicationQuestion OpportunityApplicationQuestion) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(ctx context.Context, id int64) error
}

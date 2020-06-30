package opportunities

import "github.com/joinimpact/api/internal/models"

// isPublishable returns a boolean of whether or not an opportunity (by ID) has all required fields for publishing.
func isPublishable(opportunity models.Opportunity) bool {
	if len(opportunity.Title) < 8 {
		return false
	}

	if len(opportunity.Description) < 24 {
		return false
	}

	// TODO: requirements object check

	return true
}

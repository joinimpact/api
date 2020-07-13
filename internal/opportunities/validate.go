package opportunities

import "github.com/joinimpact/api/internal/models"

// isPublishable returns a boolean of whether or not an opportunity (by ID) has all required fields for publishing. On error, it returns an array of invalid fields.
func isPublishable(opportunity models.Opportunity) ([]string, bool) {
	invalidFields := []string{}

	if len(opportunity.Title) < 8 {
		invalidFields = append(invalidFields, "title")
	}

	if len(opportunity.Description) < 24 {
		invalidFields = append(invalidFields, "description")
	}

	// Check if length of invalid fields is more than 0, and return false if it is.
	if len(invalidFields) > 0 {
		return invalidFields, false
	}

	return nil, true
}

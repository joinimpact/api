package opportunities

import "github.com/joinimpact/api/internal/models"

// isPublishable returns a boolean of whether or not an opportunity (by ID) has all required fields for publishing. On error, it returns an array of invalid fields.
func isPublishable(opportunity models.Opportunity) ([]string, bool) {
	invalidFields := []string{}

	if len(opportunity.Title) < 4 {
		invalidFields = append(invalidFields, "title")
	}

	if len(opportunity.Description) < 64 {
		invalidFields = append(invalidFields, "description")
	}

	// Check if length of invalid fields is more than 0, and return false if it is.
	if len(invalidFields) > 0 {
		return invalidFields, false
	}

	return nil, true
}

// shouldAppear returns a boolean of whether or not an opportunity should appear in internal dashboards.
func shouldAppear(view *OpportunityView) bool {
	// For now, we only want to check whether or not the opportunity has a title.
	return len(view.Title) > 1
}

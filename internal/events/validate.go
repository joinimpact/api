package events

import "github.com/joinimpact/api/internal/models"

// validateEvent validates an event and returns false if there is an issue.
func validateEvent(event *models.Event) bool {
	if event.OpportunityID == 0 {
		return false
	}

	if len(event.Title) < 4 {
		return false
	}

	return true
}

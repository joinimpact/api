package events

import "github.com/joinimpact/api/internal/models"

// requestToEvent converts a ModifyEventRequest to a models.Event.
func (s *service) requestToEvent(request ModifyEventRequest) models.Event {
	event := models.Event{}

	event.Active = true
	event.ID = request.ID
	event.OpportunityID = request.ID
	event.CreatorID = request.CreatorID
	event.Title = request.Title
	event.Description = request.Description
	event.Hours = request.Hours

	if request.EventSchedule != nil {
		event.DateOnly = request.EventSchedule.DateOnly
		event.FromDate = request.EventSchedule.FromDate
		event.ToDate = request.EventSchedule.ToDate
	}

	if request.Location != nil {
		event.LocationLatitude = request.Location.Latitude
		event.LocationLongitude = request.Location.Longitude
	}

	return event
}

func (s *service) eventToView(event models.Event) EventView {
	view := EventView{}

	return view
}

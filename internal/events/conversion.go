package events

import (
	"fmt"

	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/pkg/location"
)

// requestToEvent converts a ModifyEventRequest to a models.Event.
func (s *service) requestToEvent(request ModifyEventRequest) models.Event {
	event := models.Event{}

	event.Active = true
	event.ID = request.ID
	event.OpportunityID = request.OpportunityID
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

func (s *service) eventToView(event models.Event) (*EventView, error) {
	view := &EventView{}

	view.ID = event.ID
	view.OpportunityID = event.OpportunityID
	view.CreatorID = event.CreatorID
	view.Title = event.Title
	view.Description = event.Description
	view.Hours = event.Hours

	if event.LocationLongitude != 0 || event.LocationLatitude != 0 {
		location, err := s.locationService.CoordinatesToStreetAddress(&location.Coordinates{
			Longitude: event.LocationLongitude,
			Latitude:  event.LocationLatitude,
		})
		if err != nil {
			fmt.Println(location, err)
			return nil, NewErrServerError()
		}

		view.Location = location
	}

	view.EventSchedule = &EventSchedule{}
	view.EventSchedule.DateOnly = event.DateOnly
	view.EventSchedule.SingleDate = event.FromDate.Equal(event.ToDate)
	view.EventSchedule.FromDate = event.FromDate
	view.EventSchedule.ToDate = event.ToDate

	return view, nil
}

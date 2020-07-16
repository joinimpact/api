package location

import (
	"context"
	"errors"

	"googlemaps.github.io/maps"
)

// Service exposes methods for interacting with location data.
type Service interface {
	// CoordinatesToCity converts a Coordinates struct to a legible city, state
	// format.
	CoordinatesToCity(coordinates *Coordinates) (*Location, error)
	// CoordinatesToStreetAddress converts a Coordinates struct to a legible address, city, state
	// format.
	CoordinatesToStreetAddress(coordinates *Coordinates) (*Location, error)
}

// Options represents the options available for the Service.
type Options struct {
	APIKey string // a Google Maps API key
}

// service represents the internal implementation of the Service.
type service struct {
	client  *maps.Client
	options *Options
}

// NewService creates and returns a new Service with the provided Options struct.
func NewService(options *Options) (Service, error) {
	// Connect to the Google Maps API.
	client, err := maps.NewClient(maps.WithAPIKey(options.APIKey))
	if err != nil {
		return nil, err
	}

	return &service{
		client,
		options,
	}, nil
}

// CoordinatesToCity converts a Coordinates struct to a legible city, state
// format.
func (s *service) CoordinatesToCity(coordinates *Coordinates) (*Location, error) {
	res, err := s.client.ReverseGeocode(context.Background(), &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: coordinates.Latitude,
			Lng: coordinates.Longitude,
		},
		ResultType: []string{"locality"},
	})
	if err != nil {
		return nil, err
	}

	if len(res) < 1 {
		return nil, errors.New("could not find city")
	}

	location := Location{
		City:    &LocationName{},
		State:   &LocationName{},
		Country: &LocationName{},
	}

	for _, piece := range res[0].AddressComponents {
		switch piece.Types[0] {
		case "locality":
			location.City = &LocationName{
				ShortName: piece.ShortName,
				LongName:  piece.LongName,
			}
		case "administrative_area_level_1":
			location.State = &LocationName{
				ShortName: piece.ShortName,
				LongName:  piece.LongName,
			}
		case "country":
			location.Country = &LocationName{
				ShortName: piece.ShortName,
				LongName:  piece.LongName,
			}
		}
	}

	return &location, nil
}

// CoordinatesToStreetAddress converts a Coordinates struct to a legible address, city, state
// format.
func (s *service) CoordinatesToStreetAddress(coordinates *Coordinates) (*Location, error) {
	res, err := s.client.ReverseGeocode(context.Background(), &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: coordinates.Latitude,
			Lng: coordinates.Longitude,
		},
		ResultType: []string{"street_address"},
	})
	if err != nil {
		return nil, err
	}

	if len(res) < 1 {
		return nil, errors.New("could not find address")
	}

	location := Location{}

	for _, piece := range res[0].AddressComponents {
		switch piece.Types[0] {
		case "locality":
			location.City = &LocationName{
				ShortName: piece.ShortName,
				LongName:  piece.LongName,
			}
		case "administrative_area_level_1":
			location.State = &LocationName{
				ShortName: piece.ShortName,
				LongName:  piece.LongName,
			}
		case "country":
			location.Country = &LocationName{
				ShortName: piece.ShortName,
				LongName:  piece.LongName,
			}
		}
	}

	location.StreetAddress = &LocationName{
		ShortName: res[0].FormattedAddress,
		LongName:  res[0].FormattedAddress,
	}

	return &location, nil
}

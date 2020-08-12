package location

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

var (
	errCacheUnavailable = errors.New("cache unavailable/nil")
	errNilCoordinates   = errors.New("nil coordinates provided")
	errNilPrefix        = errors.New("nil prefix provided")
	errNilLocation      = errors.New("nil location provided")
	errItemNotFound     = errors.New("item not found in cache")
)

// DefaultCacheExpiration is the default duration of a cache item.
const DefaultCacheExpiration = 7 * 60 * 60 * 24 // number of days * 60 seconds * 60 minutes * 24 hours

// cacheSet sets an item in the cache using the prefix and Coordinates to form a key and the Location as the value.
func (s *service) cacheSet(prefix string, coordinates *Coordinates, location *Location) error {
	if s.cache == nil {
		return errCacheUnavailable
	}

	if len(prefix) < 1 {
		return errNilPrefix
	}

	if coordinates == nil || (coordinates.Latitude == 0 && coordinates.Longitude == 0) {
		return errNilCoordinates
	}

	if location == nil {
		return errNilLocation
	}

	key := cacheKey(prefix, coordinates)

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(*location)
	if err != nil {
		return err
	}

	return s.cache.Set(&memcache.Item{
		Key:        key,
		Value:      buf.Bytes(),
		Expiration: DefaultCacheExpiration,
	})
}

// cacheGet gets an item in the cache using the prefix and Coordinates to form a key.
func (s *service) cacheGet(prefix string, coordinates *Coordinates) (*Location, error) {
	if s.cache == nil {
		return nil, errCacheUnavailable
	}

	key := cacheKey(prefix, coordinates)

	item, err := s.cache.Get(key)
	if err != nil {
		return nil, errItemNotFound
	}

	var location Location
	dec := gob.NewDecoder(bytes.NewReader(item.Value))
	err = dec.Decode(&location)
	if err != nil {
		return nil, err
	}

	return &location, nil
}

// cacheKey returns a key for the cache from a prefix and Coordinates.
func cacheKey(prefix string, coordinates *Coordinates) string {
	return fmt.Sprintf("%s:%.6f,%.6f", prefix, coordinates.Latitude, coordinates.Longitude)
}

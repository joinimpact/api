package opportunities

import (
	"fmt"
	"io"
	"strings"

	"github.com/joinimpact/api/pkg/location"
)

// RecommendationQuery represents a query of opportunities for browse.
type RecommendationQuery struct {
	TagName             string
	LocationRestriction bool
	Location            *location.Coordinates
}

// buildRecommendationQuery builds an io.Reader with a json query from a Query struct.
func buildRecommendationQuery(query RecommendationQuery) io.Reader {
	filters := []string{}
	tag := ""

	if len(query.TagName) > 0 {
		tag = fmt.Sprintf(tagQuery, query.TagName)
	}

	sort := ""
	if query.Location != nil {
		sort = fmt.Sprintf(sortTemplate, query.Location.Longitude, query.Location.Latitude)
		if query.LocationRestriction {
			filters = append(filters, fmt.Sprintf(locationFilter, query.Location.Longitude, query.Location.Latitude))
		}
	}

	if len(filters) > 0 {
		// Append a trailing comma when filters apply.
		filters = append(filters, " ")
	}

	queryStr := fmt.Sprintf(recommendationQuery, tag, strings.Join(filters, ","), sort)

	fmt.Println(queryStr)

	return strings.NewReader(queryStr)
}

const recommendationQuery = `
{
	"query": {
	  "bool": {
		"must": [
		  %s
		],
		"filter": [
			%s
		  { "term": { "public": true } }
		]
	  }
	}
	%s
  }
`

const tagQuery = `
{
	"nested": {
	  "path": "tags",
	  "score_mode": "avg",
	  "query": {
		"term": {
		  "tags.name": {
			"value": "%s"
		  }
		}
	  }
	}
  }
`

const locationFilter = `
{
  "geo_distance": {
	"distance": "200km",
	"location": {
	  "lon": %f,
	  "lat": %f
	}
  }
}
`

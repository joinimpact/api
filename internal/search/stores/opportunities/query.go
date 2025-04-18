package opportunities

import (
	"fmt"
	"io"
	"strings"

	"github.com/joinimpact/api/pkg/location"
)

// Query represents a query of opportunities.
type Query struct {
	TextQuery       string                `json:"textQuery"`
	Location        *location.Coordinates `json:"location"`
	AgeRange        *AgeRange             `json:"ageRange"`
	CommitmentRange *CommitmentRange      `json:"commitmentRange"`
	Limit           uint                  `json:"-"`
	Page            uint                  `json:"-"`
}

// AgeRange represents the range of ages to filter by.
type AgeRange struct {
	Age int `json:"age"`
}

// CommitmentRange represents the range of hours to filter by.
type CommitmentRange struct {
	Minimum int `json:"min"`
	Maximum int `json:"max"`
}

// buildQuery builds an io.Reader with a json query from a Query struct.
func buildQuery(query Query) io.Reader {
	filters := []string{}
	if query.AgeRange != nil {
		filters = append(filters, fmt.Sprintf(ageFilter, query.AgeRange.Age, query.AgeRange.Age))
	}
	if query.CommitmentRange != nil {
		filters = append(filters, fmt.Sprintf(hoursFilter, query.CommitmentRange.Minimum, query.CommitmentRange.Maximum))
	}
	if len(filters) > 0 {
		// Append a trailing comma when filters apply.
		filters = append(filters, " ")
	}

	limit := uint(20)
	if query.Limit > 0 && query.Limit <= 100 {
		limit = query.Limit
	}

	limits := fmt.Sprintf(limitsTemplate, limit, query.Page*limit)

	sort := ""
	if query.Location != nil {
		sort = fmt.Sprintf(sortTemplate, query.Location.Longitude, query.Location.Latitude)
	}

	queryStr := fmt.Sprintf(queryTemplate, limits, query.TextQuery, query.TextQuery, strings.Join(filters, ","), sort)
	fmt.Println(queryStr)

	return strings.NewReader(queryStr)
}

const limitsTemplate = `
	"size": %d,
	"from": %d,
`

const queryTemplate = `
{
	%s
	"query": {
	  "bool": {
		"should": [
		  {
			"multi_match": {
			  "query": "%s",
			  "fields": ["title^4", "description^2", "organization.name^1"],
			  "zero_terms_query": "all",
			  "fuzziness": "AUTO"
			}
		  },
		  {
			"nested": {
			  "path": "tags",
			  "score_mode": "avg",
			  "query": {
				"query_string": {
				  "query": "%s",
				  "fields": ["tags.name"]
				}
			  }
			}
		  }
		],
		"filter": [
			%s
		  { "term": { "public": true } }
		],
		"minimum_should_match": 1
	  }
	}
	%s
  }
`

const sortTemplate = `
,"sort": [
	{
		"_geo_distance": {
			"location": [%f, %f],
			"order": "asc",
			"unit": "km",
			"mode": "min",
			"distance_type": "arc",
			"ignore_unmapped": true
		}
	}
]
`

const ageFilter = `
{
	"bool": {
		"should": [
			{
				"range": {
					"requirements.ageLimit.from": {
						"lte": %d
					}
				}
			},
			{
				"range": {
					"requirements.ageLimit.to": {
						"gte": %d
					}
				}
			}
		],
		"minimum_should_match": 2
	}
}
`

const hoursFilter = `
{
	"bool": {
	  "should": [
		{
		  "range": {
			"requirements.expectedHours.hours": {
			  "gte": %d
			}
		  }
		},
		{
		  "range": {
			"requirements.expectedHours.hours": {
			  "lte": %d
			}
		  }
		}
	  ],
	  "minimum_should_match": 2
	}
}
`

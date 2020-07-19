package browse

import (
	"net/http"

	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/pkg/resp"
)

// Get gets the browse page.
func Get(opportunitiesService opportunities.Service) http.HandlerFunc {
	type section struct {
		Name          string                          `json:"name"`
		Tag           string                          `json:"tag,omitempty"`
		Opportunities []opportunities.OpportunityView `json:"opportunities"`
	}
	type response struct {
		Sections []section `json:"sections"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		opportunity, err := opportunitiesService.GetOpportunity(ctx, 1283123311563247616)
		if err != nil {
			resp.OK(w, r, response{
				[]section{
					{
						Name:          "in_your_area",
						Opportunities: []opportunities.OpportunityView{},
					},
					{
						Name:          "your_interests",
						Tag:           "Research",
						Opportunities: []opportunities.OpportunityView{},
					},
				},
			})
			return
		}

		resp.OK(w, r, response{
			[]section{
				{
					Name:          "in_your_area",
					Opportunities: []opportunities.OpportunityView{*opportunity},
				},
				{
					Name:          "your_interests",
					Tag:           "Research",
					Opportunities: []opportunities.OpportunityView{*opportunity},
				},
			},
		})
	}
}

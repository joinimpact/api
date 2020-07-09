package opportunities

import (
	"net/http"

	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/parse"
	"github.com/joinimpact/api/pkg/resp"
)

// TagsPost adds tags to an organization's profile.
func TagsPost(opportunitiesService opportunities.Service) http.HandlerFunc {
	type request struct {
		Name string `json:"name" validate:"min=2,max=24"`
	}
	type response struct {
		NumberAdded int `json:"numAdded"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := struct {
			Tags []request `json:"tags"`
		}{}
		err := parse.POST(w, r, &req)
		if err != nil {
			return
		}

		opportunityID, err := idctx.Get(r, "opportunityID")
		if err != nil {
			return
		}

		tags := []string{}
		for _, tag := range req.Tags {
			tags = append(tags, tag.Name)
		}

		numAdded, err := opportunitiesService.AddOpportunityTags(ctx, opportunityID, tags)
		if err != nil {
			switch err.(type) {
			case *opportunities.ErrOpportunityNotFound:
				resp.NotFound(w, r, resp.Error(404, err.Error()))
			case *opportunities.ErrServerError:
				resp.ServerError(w, r, resp.Error(500, err.Error()))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, response{numAdded})
	}
}

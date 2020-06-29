package tags

import (
	"net/http"
	"strconv"

	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/tags"
	"github.com/joinimpact/api/pkg/resp"
)

type getTagsResponse struct {
	Tags []models.Tag `json:"tags"`
}

// GetTags creates a new organization.
func GetTags(tagsService tags.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		limit := 10
		limitString := r.URL.Query().Get("limit")
		if limitString != "" {
			limInt, err := strconv.ParseInt(limitString, 10, 8)
			if err != nil {
				resp.ServerError(w, r, resp.UnknownError)
				return
			}

			limit = int(limInt)
		}

		tags, err := tagsService.GetTags(query, limit)
		if err != nil {
			resp.ServerError(w, r, resp.UnknownError)
			return
		}

		resp.OK(w, r, &getTagsResponse{tags})
	}
}

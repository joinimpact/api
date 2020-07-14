package db

import (
	"math"
	"net/http"
	"strconv"

	"github.com/joinimpact/api/pkg/dbctx"
	"github.com/joinimpact/api/pkg/resp"
)

// ContextMiddleware converts get parameters to dbcontext Requests.
func ContextMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Set default limit.
			limit := 20
			// Get the limit from the query.
			limitString := r.URL.Query().Get("limit")
			if limitString != "" {
				limInt, err := strconv.ParseInt(limitString, 10, 8)
				if err != nil {
					resp.BadRequest(w, r, resp.Error(400, "invalid limit parameter, must be an integer"))
					return
				}

				// Clamp the limit between 1 and 100.
				limit = int(math.Min(math.Max(float64(limInt), 1), 100))
			}

			// Set default page.
			page := 0
			// Get the page from the query.
			pageString := r.URL.Query().Get("page")
			if pageString != "" {
				pageInt, err := strconv.ParseInt(pageString, 10, 8)
				if err != nil {
					resp.BadRequest(w, r, resp.Error(400, "invalid page parameter, must be an integer"))
					return
				}

				page = int(pageInt)
			}

			// Get the query parameter.
			queryString := r.URL.Query().Get("query")

			// Inject the limit value.
			ctx = dbctx.Inject(ctx, dbctx.Request{
				Limit: limit,
				Page:  page,
				Query: queryString,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

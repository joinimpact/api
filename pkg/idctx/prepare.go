package idctx

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/joinimpact/api/pkg/resp"
)

type key string

const (
	keyIDContext key = "idctx-main-context"
)

type idMap map[string]int64

// Prepare is an HTTP middleware that processes IDs in URL params.
func Prepare(names ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Get the idMap from the context.
			ids, ok := ctx.Value(keyIDContext).(idMap)
			if !ok {
				// If the idMap doesn't exist, create it.
				ids = idMap{}
			}

			for _, name := range names {
				stringParam := chi.URLParam(r, name)

				// Convert the string parameter to an int64.
				convInt, err := strconv.ParseInt(stringParam, 10, 64)
				if err != nil {
					resp.BadRequest(w, r, resp.Error(400, fmt.Sprintf("invalid %s", name)))
					return
				}

				ids[name] = convInt
			}

			// Add the idMap to the context and serve the next handler.
			ctx = context.WithValue(ctx, keyIDContext, ids)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

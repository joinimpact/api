package idctx

import (
	"context"
	"errors"
	"net/http"
)

// Get gets an ID by key from a Request and returns it.
func Get(r *http.Request, key string) (int64, error) {
	ctx := r.Context()

	return GetFromContext(ctx, key)
}

// GetFromContext gets an ID by key from a Context and returns it.
func GetFromContext(ctx context.Context, key string) (int64, error) {
	// Get the idMap from the context.
	ids, ok := ctx.Value(keyIDContext).(idMap)
	if !ok {
		// If the idMap doesn't exist, return an error.
		return 0, errors.New("invalid context")
	}

	return ids[key], nil
}

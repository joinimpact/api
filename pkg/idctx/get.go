package idctx

import (
	"errors"
	"net/http"
)

// Get gets an ID by key from a Request and returns it.
func Get(r *http.Request, key string) (int64, error) {
	ctx := r.Context()

	// Get the idMap from the context.
	ids, ok := ctx.Value(keyIDContext).(idMap)
	if !ok {
		// If the idMap doesn't exist, return an error.
		return 0, errors.New("invalid context")
	}

	return ids[key], nil
}

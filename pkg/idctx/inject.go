package idctx

import "context"

// Inject allows custom middleware to set an ID in the idctx context with a
// key and value, and returns a modified context with the ID injected.
func Inject(ctx context.Context, key string, value int64) context.Context {
	// Get the idMap from the context.
	ids, ok := ctx.Value(keyIDContext).(idMap)
	if !ok {
		// If the idMap doesn't exist, create it.
		ids = idMap{}
	}

	ids[key] = value

	return context.WithValue(ctx, keyIDContext, ids)
}

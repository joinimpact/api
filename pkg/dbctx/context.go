package dbctx

import "context"

type key int

const (
	keyDBContext key = iota
)

// Defaults
const (
	DefaultLimit = 20
	DefaultPage  = 0
	DefaultQuery = ""
)

// Request provides different options for returning database results.
type Request struct {
	Limit int
	Page  int
	Query string
}

// Inject adds a dbctx.Request to a context and returns it.
func Inject(ctx context.Context, request Request) context.Context {
	return context.WithValue(ctx, keyDBContext, &request)
}

// Get gets a dbctx.Request from a context.Context, and returns a default Request on failure.
func Get(ctx context.Context) *Request {
	value, ok := ctx.Value(keyDBContext).(*Request)
	if !ok {
		return &Request{
			DefaultLimit,
			DefaultPage,
			DefaultQuery,
		}
	}

	return value
}

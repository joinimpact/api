package pubsub

import "context"

// EventName represents the identifier of a single event.
// It can be a single uppercase word with underscores, such as NEW_MESSAGE,
// and can also contain a category identifier, ex: messages.NEW_MESSAGE
type EventName string

// Event represents a single event passing through the pubsub system.
type Event struct {
	EventName EventName   `json:"event"`
	Payload   interface{} `json:"payload"`
	ctx       context.Context
}

// Context returns the event's context.
func (e *Event) Context() context.Context {
	if e.ctx != nil {
		return e.ctx
	}

	return context.Background()
}

// WithContext returns a shallow copy of e with the context set to the provided ctx.
// The provided ctx must be non-nil, or WithContext will panic.
func (e *Event) WithContext(ctx context.Context) *Event {
	if ctx == nil {
		panic("nil context")
	}
	e2 := new(Event)
	*e2 = *e
	e2.ctx = ctx
	return e2
}

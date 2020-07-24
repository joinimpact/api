package pubsub

import (
	"errors"
	"sync"
)

// Errors
var (
	ErrBrokerClosed = errors.New("broker closed")
)

// Stream represents a collection of related events.
type Stream string

// Broker represents an interface for sending and receiving messages through
// the pubsub system.
type Broker interface {
	// Subscribe subscribes to an event stream and returns a channel of Events.
	Subscribe(stream Stream) <-chan Event
	// Publish publishes an event to the specified stream.
	Publish(stream Stream, event Event) error
	// Close gracefully closes the broker and all channels.
	Close()
}

// broker represents the internal implementation of the Broker interface.
type broker struct {
	mu       sync.RWMutex
	registry map[Stream][]chan Event
	closed   bool
}

// NewBroker creates and returns a new Broker.
func NewBroker() Broker {
	return &broker{
		sync.RWMutex{},
		make(map[Stream][]chan Event),
		false,
	}
}

// Subscribe subscribes to an event stream and returns a channel of Events.
func (b *broker) Subscribe(stream Stream) <-chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	channel := make(chan Event)

	b.registry[stream] = append(b.registry[stream], channel)

	return channel
}

// Publish publishes an event to the specified stream.
func (b *broker) Publish(stream Stream, event Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return ErrBrokerClosed
	}

	for _, ch := range b.registry[stream] {
		go func(ch chan Event) {
			ch <- event
		}(ch)
	}

	return nil
}

// Close gracefully closes the broker and all channels.
func (b *broker) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.closed {
		b.closed = true
		for _, subs := range b.registry {
			for _, ch := range subs {
				close(ch)
			}
		}
	}
}

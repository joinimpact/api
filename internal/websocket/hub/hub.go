package hub

import (
	"errors"

	"github.com/joinimpact/api/internal/websocket"
)

// Options contains optional parameters for configuring the hub
type Options struct {
}

// Hub represents a map of channels containing subscriptions
type Hub struct {
	channels map[ChannelID]*Channel
	options  Options
}

// NewHub creates and returns a new Hub with the provided Options object
func NewHub(options Options) *Hub {
	return &Hub{
		channels: make(map[ChannelID]*Channel),
		options:  options,
	}
}

// newChannel creates and returns a channel in the Hub
func (h *Hub) newChannel(id ChannelID) (*Channel, error) {
	if c, ok := h.channels[id]; ok || c != nil {
		return nil, errors.New("channel already exists")
	}

	c := InitChannel(id)
	h.channels[id] = c

	return c, nil
}

// getChannel returns a Channel by ID, or an error if it is not found
func (h *Hub) getChannel(id ChannelID) (*Channel, error) {
	c, ok := h.channels[id]
	if !ok || c == nil {
		return h.newChannel(id)
	}

	return c, nil
}

// removeChannel removes a Channel from the Hub
func (h *Hub) removeChannel(id ChannelID) error {
	delete(h.channels, id)
	return nil
}

// Subscribe subscribes a client to a Channel
func (h *Hub) Subscribe(id ChannelID, session *Session) error {
	c, err := h.getChannel(id)
	if err != nil {
		return err
	}

	// Subscribe the session to the channel
	c.Subscribe(session)

	return nil
}

// Unsubscribe unsubscribes a client from a Channel
func (h *Hub) Unsubscribe(id ChannelID, sessionID SessionID) error {
	c, err := h.getChannel(id)
	if err != nil {
		return err
	}

	// Unsubscribe the session from the Channel
	empty := c.Unsubscribe(sessionID)
	if empty {
		// If the Channel is empty, we'll remove it from the Hub
		err := h.removeChannel(id)
		if err != nil {
			return err
		}
	}

	return nil
}

// RouteMessage routes a message to the intended Channel
func (h *Hub) RouteMessage(id ChannelID, message websocket.Message) error {
	c, err := h.getChannel(id)
	if err != nil {
		return err
	}

	// Send the message to the channel
	c.In <- message
	return nil
}

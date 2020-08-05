package hub

import (
	"sync"

	"github.com/joinimpact/api/internal/websocket"
)

// ChannelID is a string ID reference to one channel
type ChannelID string

// Channel wraps a group of subscriptions
type Channel struct {
	ChannelID     ChannelID
	In            chan websocket.Message
	lock          *sync.RWMutex
	subscriptions map[SessionID]*Session
}

// Types
const (
	TypeMessage = "MESSAGE"
)

// InitChannel starts and returns a Channel
// This function starts the Reader in a go routine as well,
// which will begin to watch the channel and send messages
// to connected sessions
func InitChannel(id ChannelID) *Channel {
	c := &Channel{
		ChannelID:     id,
		In:            make(chan websocket.Message, 128),
		lock:          &sync.RWMutex{},
		subscriptions: make(map[SessionID]*Session),
	}

	// Launch the Reader
	go c.Reader()

	return c
}

// Fanout fans a message out to all connected sessions in the Channel
func (c *Channel) Fanout(m websocket.Message) {
	// Loop through all connected sessions
	for _, s := range c.subscriptions {
		if s.Closed {
			continue
		}
		// Send message to session
		s.SendMessage(m)
	}
}

// Reader watches the In channel and routes the messages to the connected sessions
func (c *Channel) Reader() {
	for {
		// Get the message from the input channel
		m := <-c.In

		c.Fanout(m)
	}
}

// Subscribe subscribes a session to the channel
func (c *Channel) Subscribe(s *Session) {
	// Lock the mutex, blocking writes
	c.lock.Lock()
	// Unlock the mutex after the function completes
	defer c.lock.Unlock()
	// Add the session to the subscriptions list
	c.subscriptions[s.SessionID] = s
}

// Unsubscribe unsubscribes a session from the channel
// The function returns true if the channel has no sessions left
// so that the hub can then delete the channel from the master map
func (c *Channel) Unsubscribe(id SessionID) bool {
	// Lock the mutex, blocking writes
	c.lock.Lock()
	// Unlock the mutex after the function completes
	defer c.lock.Unlock()

	s, ok := c.subscriptions[id]
	if !ok {
		// Explanation of return at bottom of function
		return len(c.subscriptions) == 0
	}

	// Update the session's subscribed Channels list
	s.Susbcriptions[c.ChannelID] = false

	// Remove the session from the subscriptions list by ID
	delete(c.subscriptions, id)

	// Return true if the number of subscriptions is 0
	// This is done so that the hub can remove channels
	// that have no active sessions to serve
	return len(c.subscriptions) == 0
}

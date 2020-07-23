package hubmanager

import (
	"errors"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/joinimpact/api/internal/websocket/hub"
)

// userID represents an ID of a single User
type userID int64

// HubManager manages users and their subscriptions
type HubManager struct {
	hub      *hub.Hub
	lock     *sync.RWMutex
	sessions map[userID]map[hub.SessionID]*hub.Session // stores Sessions by user IDs
}

// NewHubManager creates and returns a new HubManager with the hub provided
func NewHubManager(h *hub.Hub) *HubManager {
	return &HubManager{
		h,
		&sync.RWMutex{},
		make(map[userID]map[hub.SessionID]*hub.Session), // initialize a blank map
	}
}

// Register registers a Session under a userID
func (h *HubManager) Register(id int64, session *hub.Session) error {
	// Lock the mutex, blocking writes
	h.lock.Lock()
	// Unlock the mutex after the function completes
	defer h.lock.Unlock()
	uid := userID(id)
	m, ok := h.sessions[uid]
	if !ok || m == nil {
		// If the user's ID is not registered, create a new registry for
		// the user
		h.sessions[uid] = make(map[hub.SessionID]*hub.Session)
	}

	// Add the Session to the registry
	h.sessions[uid][session.SessionID] = session

	return nil
}

// Unregister unregisters a Session from a userID group
func (h *HubManager) Unregister(id int64, sessionID hub.SessionID) error {
	// Lock the mutex, blocking writes
	h.lock.Lock()
	// Unlock the mutex after the function completes
	defer h.lock.Unlock()
	uid := userID(id)
	m, ok := h.sessions[uid]
	if !ok || m == nil {
		return errors.New("group does not exist")
	}

	// Remove the Session from the group
	delete(h.sessions[uid], sessionID)

	if len(h.sessions[uid]) <= 0 {
		// If the group is empty, delete it
		delete(h.sessions, uid)
	}

	return nil
}

// subscribeAll subscribes all Sessions in a user ID group to a Channel
func (h *HubManager) subscribeAll(id userID, channelID hub.ChannelID) {
	for _, v := range h.sessions[id] {
		// Subscribe the Session to the Hub Channel
		h.hub.Subscribe(channelID, v)
	}
}

// unsubscribeAll unsubscribes all Sessions in a user ID group from a Channel
func (h *HubManager) unsubscribeAll(id userID, channelID hub.ChannelID) {
	for _, v := range h.sessions[id] {
		// Unsubscribe the Session from the Hub Channel
		h.hub.Unsubscribe(channelID, v.SessionID)
	}
}

// processMessage modifies subscriptions based on pubsub.Messages received
func (h *HubManager) processMessage(m pubsub.Message) {
	return
	// Decode the event code
	// _, err := strconv.ParseInt(m.Header["event"], 10, 64)
	// if err != nil {
	// codeDecoded = -1
	// }

	// Convert from an int64 to an int
	// code := int(codeDecoded)

	// switch code {
	// case codes.EventAddedToConversation:
	// 	d, err := conversationmember.Decode(m)
	// 	if err != nil {
	// 		// Log error if the conversation member is unable to be decoded
	// 		log.Printf("[HubManager] Error while decoding conversation member: %e", err)
	// 		return
	// 	}

	// 	groupID := userID(d.UserId)

	// 	// Subscribe the Sessions to the Channel
	// 	h.subscribeAll(groupID, hub.ChannelID(d.ConversationId))
	// case codes.EventRemovedFromConversation:
	// 	d, err := conversationmember.Decode(m)
	// 	if err != nil {
	// 		// Log error if the conversation member is unable to be decoded
	// 		log.Printf("[HubManager] Error while decoding conversation member: %e", err)
	// 		return
	// 	}

	// 	groupID := userID(d.UserId)

	// 	// Unsubscribe the Sessions from the Channel
	// 	h.unsubscribeAll(groupID, hub.ChannelID(d.ConversationId))
	// }
}

// messagePump takes messages in from a provided channel, interprets the
// results, and forwards the messages to the Hub
func (h *HubManager) messagePump(channel chan pubsub.Message) {
	for {
		select {
		case m := <-channel:
			h.processMessage(m)

			// to, ok := m.Header["to"]
			// if !ok {
			// 	// If the routing header is not present, ignore the message
			// 	continue
			// }

			// channelID, err := strconv.ParseInt(to, 10, 64)
			// if err != nil {
			// 	// The channelID is invalid, so ignore the message
			// 	continue
			// }
			// // Pass the message into the Hub
			// h.hub.RouteMessage(hub.ChannelID(channelID), m)

			// if toUser, ok := m.Header["to_user"]; ok {
			// 	// If we also receive a to_user to send to a specific user
			// 	channelID, err := strconv.ParseInt(toUser, 10, 64)
			// 	if err != nil {
			// 		// The channelID is invalid, so don't send
			// 		continue
			// 	}
			// 	// Pass the message into the Hub
			// 	h.hub.RouteMessage(hub.ChannelID(channelID), m)
			// }
		}
	}
}

// StartMessagePump starts pumping messages into the Hub from the provided
// pubsub Message channel
func (h *HubManager) StartMessagePump(channel chan pubsub.Message) {
	// Start the messagePump in a goroutine
	go h.messagePump(channel)
}

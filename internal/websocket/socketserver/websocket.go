package socketserver

import (
	gws "github.com/gorilla/websocket"
	"github.com/joinimpact/api/internal/authentication"
	"github.com/joinimpact/api/internal/conversations"
	"github.com/joinimpact/api/internal/organizations"
	"github.com/joinimpact/api/internal/pubsub"
	"github.com/joinimpact/api/internal/websocket"
	"github.com/joinimpact/api/internal/websocket/hub"
	"github.com/joinimpact/api/internal/websocket/hubmanager"
)

// Timeouts etc.
const (
	heartbeatTimeout = 10000
)

// WebSocket error codes
const (
	ErrUnsupported         = 1003
	ErrInvalidAuthHeader   = 4001
	ErrUnauthorized        = 4002
	ErrUnableToGetChannels = 4003
	ErrTimeout             = 4004
)

var stream pubsub.Stream = pubsub.Stream("impact.users")

// WebSocketManager provides methods for reading and writing from WebSockets
type WebSocketManager struct {
	hub                   *hub.Hub
	hubManager            *hubmanager.HubManager
	broker                pubsub.Broker
	authenticationService authentication.Service
	organizationsService  organizations.Service
	conversationsService  conversations.Service
}

// NewWebSocketManager creates and returns a WebSocketManager based on the
// provided dependencies
func NewWebSocketManager(h *hub.Hub, hm *hubmanager.HubManager, broker pubsub.Broker, authenticationService authentication.Service, organizationsService organizations.Service, conversationsService conversations.Service) *WebSocketManager {
	return &WebSocketManager{
		h,
		hm,
		broker,
		authenticationService,
		organizationsService,
		conversationsService,
	}
}

// SubscribeHub subscribes the hub to the pubsub broker.
func (w *WebSocketManager) SubscribeHub() {
	channel := w.broker.Subscribe(stream)

	w.hubManager.StartMessagePump(channel)
}

// Reader reads and parses messages from the WebSocket
func (w *WebSocketManager) Reader(s *hub.Session) {
	conn := s.Conn

	var m websocket.Message
	for {
		if s.Closed {
			// Return if the session was closed.
			return
		}

		// Attempt to read the client message as JSON
		err := conn.ReadJSON(&m)
		if err != nil {
			// If the client data is not valid JSON, close the connection
			// with an ErrUnsupported
			s.CloseWithError(ErrUnsupported)
			return
		}

		switch m.Opcode {
		case websocket.OpcodeHeartbeat:
			// Reset the timer
			s.Timer.Reset(s.HeartbeatTimeout)
			// Send heartbeat acknowledge message
			s.SendMessage(websocket.Message{
				Opcode: websocket.OpcodeHeartbeatAck,
			})
		case websocket.OpcodeClientAuthenticate:
			data := m.Data.(map[string]interface{})
			w.WebSocketAuthenticate(s, data)
		}
	}
}

// Writer writes messages to the WebSocket
func (w *WebSocketManager) Writer(s *hub.Session) {
	// Send the message to the client
	s.SendMessage(websocket.HelloMessage)

	for {
		select {
		case m := <-s.Channel:
			// Write the message
			s.Conn.WriteJSON(m)
		case e := <-s.Close:
			// Write the message to the WebSocket connection
			s.Conn.WriteMessage(gws.CloseMessage, gws.FormatCloseMessage(e, ""))

			// Close the connection
			s.Conn.Close()

			// Close the channels
			close(s.Close)
			close(s.Channel)

			if s.UserID > 0 {
				// Unsubscribe from the HubManager if the user is authenticated
				w.hubManager.Unregister(s.UserID, s.SessionID)
			}

			// Unsubscribe from all Channels
			for id, active := range s.Susbcriptions {
				if !active {
					// If the subscription is already inactive, ignore
					continue
				}

				// Unsubscribe the Session from the Channel
				w.hub.Unsubscribe(id, s.SessionID)
			}

			// Return the function
			return
		case <-s.Timer.C:
			// The heartbeat timeout has expired, send an error and close the
			// WebSocket connection
			s.CloseWithError(ErrTimeout)
		}
	}
}

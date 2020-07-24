package socketserver

import (
	"fmt"

	"github.com/joinimpact/api/internal/websocket"
	"github.com/joinimpact/api/internal/websocket/hub"
)

// WebSocketHandler is an interface for creating handlers for WebSockets
type WebSocketHandler func(s *hub.Session, data map[string]string) error

// WebSocketAuthenticate attempts to authenticate a WebSocket connection
func (w *WebSocketManager) WebSocketAuthenticate(s *hub.Session, data map[string]interface{}) error {
	// Authenticate the token with the Authenticator service, and return an
	// error on failure
	userID, err := w.authenticationService.GetUserIDFromToken(data["token"].(string))
	if err != nil {
		s.CloseWithError(ErrUnauthorized)
		return nil
	}

	// Set the Session's UserID to the UserId in the claims
	s.UserID = userID

	// // Get all of a user's channels with the UserID in the token claims
	channels, err := w.getUserChannels(userID)
	if err != nil {
		s.CloseWithError(ErrUnableToGetChannels)
		return nil
	}

	// Subscribe the Session to the HubManager group for their user ID
	if err := w.hubManager.Register(userID, s); err != nil {
		s.CloseWithError(ErrUnableToGetChannels)
		return err
	}

	// Subscribe the Session to their own Channel
	w.hub.Subscribe(hub.ChannelID(fmt.Sprintf("user/%d", userID)), s)

	for _, id := range channels {
		// Subscribe the Session to each Channel returned
		w.hub.Subscribe(id, s)
	}

	// Create a message to send to the client
	m := websocket.Message{
		Opcode: websocket.OpcodeAuthenticationSuccess,
		Data: map[string]int64{
			"userId": userID,
		},
	}

	// Send the message
	s.SendMessage(m)

	return nil
}

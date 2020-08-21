package hub

import (
	"time"

	gsw "github.com/gorilla/websocket"
	"github.com/joinimpact/api/internal/websocket"
	"github.com/joinimpact/api/pkg/scopes"
	"github.com/rs/zerolog/log"
)

// SessionID is an int64 ID reference to one session/user
type SessionID int64

// Session represents a connected user
type Session struct {
	SessionID        SessionID          // the unique ID of the session
	UserID           int64              // the user ID of the user associated with the session
	Conn             *gsw.Conn          // the actual WebSocket connection
	Closed           bool               // whether or not the session was closed
	Close            chan int           // channel for closing the WebSocket connection with an int error
	Channel          chan interface{}   // channel for sending messages to the WebSocket
	SequenceNumber   int                // allows the client to know when messages are skipped
	HeartbeatTimeout time.Duration      // timeout for heartbeating
	Timer            *time.Timer        // timer for heartbeating
	Susbcriptions    map[ChannelID]bool // list of subscribed channels
}

// SendMessage sends a message to the websocket client
func (s *Session) SendMessage(message websocket.Message) error {
	message.SequenceNumber = s.SequenceNumber
	s.SequenceNumber++

	marshaled := scopes.Marshal(scopes.ScopeAuthenticated, message)

	log.Logger.Info().Fields(map[string]interface{}{
		"userId":        s.UserID,
		"subscriptions": s.Susbcriptions,
		"payload":       message,
	}).Msg("Debug: message sent to WebSocket session")

	// Push the message to the channel
	s.Channel <- marshaled
	return nil
}

// CloseWithError sends an error to the websocket client and closes the connection
func (s *Session) CloseWithError(errorCode int) error {
	if s.Closed {
		return nil
	}

	// Push the message to the closing channel
	s.Close <- errorCode
	s.Closed = true
	return nil
}

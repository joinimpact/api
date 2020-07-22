package hub

import "github.com/gorilla/websocket"

// SessionID is an int64 ID reference to one session/user
type SessionID int64

// Session represents a connected user
type Session struct {
	SessionID      SessionID
	Conn           *websocket.Conn
	SequenceNumber int
}

// SendMessage sends a message to the websocket client
func (s *Session) SendMessage(message interface{}) error {
	return s.Conn.WriteJSON(message)
}

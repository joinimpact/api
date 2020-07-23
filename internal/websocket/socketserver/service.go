package socketserver

import "net/http"

// Service represents a websocket service.
type Service interface {
	// Handler is the main entrypoint for the WebSocket.
	Handler() http.HandlerFunc
}

// service represents the internal implementation of the socketserver.Service interface.
type service struct {
	wsm *WebSocketManager
}

// NewService creates and returns a new Service which wraps the provdied WebSocketManager.
func NewService(wsm *WebSocketManager) Service {
	return &service{wsm}
}

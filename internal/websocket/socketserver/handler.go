package socketserver

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	gws "github.com/gorilla/websocket"
	"github.com/joinimpact/api/internal/websocket"
	"github.com/joinimpact/api/internal/websocket/hub"
)

// upgrader allows for upgrading HTTP requests to WebSocket connections
var upgrader = gws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// getToken attempts to get the token from the Authorization HTTP header.
func getToken(r *http.Request) (string, error) {
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Bearer" {
		return "", errors.New("can not get header")
	}

	return auth[1], nil
}

// Handler is the main entrypoint for the WebSocket.
func (s *service) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Attempt to upgrade the connection to a WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		// Generate a new session ID
		sessionID := hub.SessionID(uuid.New().ID())

		// Calculate the heartbeat timeout
		d := time.Duration(websocket.HeartbeatTimeout) * time.Millisecond

		d = d + 1000*time.Millisecond

		// Create a hub.Session object with the session ID and WebSocket connection
		session := hub.Session{
			SessionID:        sessionID,
			UserID:           0,
			Conn:             conn,
			Closed:           false,
			Channel:          make(chan interface{}, 4),
			Close:            make(chan int, 4),
			HeartbeatTimeout: d,
			Timer:            time.NewTimer(d),
			Susbcriptions:    make(map[hub.ChannelID]bool),
		}

		token, err := getToken(r)
		if err == nil {
			s.wsm.WebSocketAuthenticate(&session, map[string]interface{}{
				"token": token,
			})
		}

		// Launch the Reader and Writer asynchronously
		go s.wsm.Reader(&session)
		go s.wsm.Writer(&session)
	}
}

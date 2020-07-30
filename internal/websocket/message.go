package websocket

// Opcodes
const (
	OpcodeHello                 = iota
	OpcodeClientAuthenticate    = iota
	OpcodeAuthenticationSuccess = iota
	OpcodeHeartbeat             = iota
	OpcodeHeartbeatAck          = iota
	OpcodeEvent                 = iota
)

// HelloMessage is the message sent when the client first connects to the server.
var HelloMessage Message = Message{
	Opcode: OpcodeHello,
	Data: map[string]interface{}{
		"heartbeatInterval": HeartbeatTimeout,
		"instructions":      "Welcome! The server will begin to expect a Heartbeat (4) opcode every heartbeatInterval miliseconds.",
	},
}

// Message represents a WebSocket message.
type Message struct {
	Opcode         int         `json:"op"`
	EventName      string      `json:"e,omitempty"`
	Data           interface{} `json:"d"`
	SequenceNumber int         `json:"sequenceNumber"`
}

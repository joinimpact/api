package websocket

// Message represents a WebSocket message.
type Message struct {
	Opcode         int         `json:"op"`
	EventName      string      `json:"e,omitempty"`
	Data           interface{} `json:"d"`
	SequenceNumber int         `json:"sequenceNumber"`
}

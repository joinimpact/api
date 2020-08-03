package conversations

import "time"

// MessageView represents a view of a message.
type MessageView struct {
	ID                int64                  `json:"id"`
	ConversationID    int64                  `json:"conversationId"`
	SenderID          int64                  `json:"senderId"`
	SenderPerspective uint                   `json:"senderPerspective"`
	Timestamp         time.Time              `json:"timestamp"`
	Type              string                 `json:"type"`
	Edited            bool                   `json:"edited"`
	EditedTimestamp   time.Time              `json:"editedTimestamp"`
	Body              map[string]interface{} `json:"body"`
}

package hubmanager

import (
	"fmt"

	"github.com/joinimpact/api/internal/websocket/hub"
)

// ConversationIDToChannelID converts a conversation ID to a ChannelID.
func ConversationIDToChannelID(id int64) hub.ChannelID {
	return hub.ChannelID(fmt.Sprintf("conversation/%d", id))
}

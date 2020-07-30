package socketserver

import (
	"github.com/joinimpact/api/internal/websocket/hub"
	"github.com/joinimpact/api/internal/websocket/hubmanager"
)

// getUserChannels gets a user's channels.
func (wsm *WebSocketManager) getUserChannels(userID int64) ([]hub.ChannelID, error) {
	ids := []hub.ChannelID{}

	memberships, err := wsm.conversationsService.GetUserConversationMemberships(userID)
	if err != nil {
		return nil, err
	}

	for _, membership := range memberships {
		ids = append(ids, hubmanager.ConversationIDToChannelID(membership.ConversationID))
	}

	organizations, err := wsm.organizationsService.GetUserOrganizations(userID)
	if err != nil {
		return nil, err
	}

	for _, organization := range organizations {
		conversations, err := wsm.conversationsService.GetOrganizationConversations(organization.ID)
		if err != nil {
			return nil, err
		}

		for _, conversation := range conversations {
			ids = append(ids, hubmanager.ConversationIDToChannelID(conversation.ID))
		}
	}

	return ids, nil
}

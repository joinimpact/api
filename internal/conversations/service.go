package conversations

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joinimpact/api/internal/cdn"
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/email"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/pubsub"
	"github.com/joinimpact/api/internal/snowflakes"
	"github.com/joinimpact/api/pkg/dbctx"
	"github.com/joinimpact/api/pkg/location"
	"github.com/rs/zerolog"
)

var stream = pubsub.Stream("impact.users")

// Events
const (
	MessageSent = "messages.MESSAGE_SENT"
)

// Service defines methods for interacting with conversations and messages.
type Service interface {
	// CreateOpportunityMembershipRequestConversation creates an opportunity membership request conversation and adds a message to it. Returns conversation ID on success.
	CreateOpportunityMembershipRequestConversation(ctx context.Context, organizationID, opportunityID, opportunityMembershipRequestID, volunteerID int64, messageStr string) (int64, error)
	// GetUserConversationMemberships gets a user's volunteer conversation memberships.
	GetUserConversationMemberships(userID int64) ([]models.ConversationMembership, error)
	// GetUserConversations gets all of a user's conversations.
	GetUserConversations(ctx context.Context, userID int64) (*ConversationsResponse, error)
	// GetUserConversation gets a single conversation from a user perspective.
	GetUserConversation(ctx context.Context, conversationID int64) (*ConversationView, error)
	// GetOrganizationConversations gets an organization's internal conversations.
	GetOrganizationConversations(ctx context.Context, organizationID int64) (*ConversationsResponse, error)
	// GetOrganizationConversation gets a single conversation from a conversation perspective.
	GetOrganizationConversation(ctx context.Context, conversationID int64) (*ConversationView, error)
	// SendStandardMessage sends a standard message to a conversation, returning the ID on success.
	SendStandardMessage(ctx context.Context, conversationID, senderID int64, messageText string, asOrganization bool) (int64, error)
	// GetConversationMessages gets messages by conversation ID.
	GetConversationMessages(ctx context.Context, conversationID int64) (*ConversationMessagesResponse, error)
	// SendHoursRequestMessage sends an hours request message to a user's organization message.
	SendHoursRequestMessage(ctx context.Context, userID, organizationID, requestID int64) (int64, error)
	// SendHoursRequestAcceptedMessage sends an hours request accept message to a user's organization message.
	SendHoursRequestAcceptedMessage(ctx context.Context, userID, requestID int64) (int64, error)
	// SendHoursRequestDeclinedMessage sends an hours request decline message to a user's organization message.
	SendHoursRequestDeclinedMessage(ctx context.Context, userID, requestID int64) (int64, error)
}

// service represents the internal implementation of the conversations Service.
type service struct {
	conversationRepository                             models.ConversationRepository
	conversationMembershipRepository                   models.ConversationMembershipRepository
	conversationOpportunityMembershipRequestRepository models.ConversationOpportunityMembershipRequestRepository
	conversationOrganizationMembershipRepository       models.ConversationOrganizationMembershipRepository
	messageRepository                                  models.MessageRepository
	opportunityRepository                              models.OpportunityRepository
	userRepository                                     models.UserRepository
	userProfileFieldRepository                         models.UserProfileFieldRepository
	userTagRepository                                  models.UserTagRepository
	tagRepository                                      models.TagRepository
	volunteeringHourLogRequestRepository               models.VolunteeringHourLogRequestRepository
	config                                             *config.Config
	logger                                             *zerolog.Logger
	snowflakeService                                   snowflakes.SnowflakeService
	emailService                                       email.Service
	broker                                             pubsub.Broker
	locationService                                    location.Service
	cdnClient                                          *cdn.Client
}

// NewService creates and returns a new conversations.Service.
func NewService(conversationRepository models.ConversationRepository, conversationMembershipRepository models.ConversationMembershipRepository, conversationOpportunityMembershipRequestRepository models.ConversationOpportunityMembershipRequestRepository, conversationOrganizationMembershipRepository models.ConversationOrganizationMembershipRepository, messageRepository models.MessageRepository, opportunityRepository models.OpportunityRepository, userRepository models.UserRepository, userProfileFieldRepository models.UserProfileFieldRepository, userTagRepository models.UserTagRepository, tagRepository models.TagRepository, volunteeringHourLogRequestRepository models.VolunteeringHourLogRequestRepository, config *config.Config, logger *zerolog.Logger, snowflakeService snowflakes.SnowflakeService, emailService email.Service, broker pubsub.Broker, locationService location.Service) Service {
	return &service{
		conversationRepository,
		conversationMembershipRepository,
		conversationOpportunityMembershipRequestRepository,
		conversationOrganizationMembershipRepository,
		messageRepository,
		opportunityRepository,
		userRepository,
		userProfileFieldRepository,
		userTagRepository,
		tagRepository,
		volunteeringHourLogRequestRepository,
		config,
		logger,
		snowflakeService,
		emailService,
		broker,
		locationService,
		cdn.NewCDNClient(config),
	}
}

// CreateOpportunityMembershipRequestConversation creates an opportunity membership request conversation and adds a message to it. Returns conversation ID on success.
func (s *service) CreateOpportunityMembershipRequestConversation(ctx context.Context, organizationID, opportunityID, opportunityMembershipRequestID, volunteerID int64, messageStr string) (int64, error) {
	// Conversation
	conversation := models.Conversation{}
	conversation.ID = s.snowflakeService.GenerateID()
	conversation.Active = true
	conversation.OrganizationID = organizationID
	conversation.Type = 1 // TODO: make this a constant in the models package

	if err := s.conversationRepository.Create(conversation); err != nil {
		return 0, NewErrServerError()
	}

	// ConversationMembership
	conversationMembership := models.ConversationMembership{}
	conversationMembership.ID = s.snowflakeService.GenerateID()
	conversationMembership.Active = true
	conversationMembership.ConversationID = conversation.ID
	conversationMembership.UserID = volunteerID
	conversationMembership.Role = 0 // TODO: make this a constant in the models package

	if err := s.conversationMembershipRepository.Create(conversationMembership); err != nil {
		return 0, NewErrServerError()
	}

	// ConversationOpportunityMembershipRequest
	conversationOpportunityMembershipRequest := models.ConversationOpportunityMembershipRequest{}
	conversationOpportunityMembershipRequest.ID = s.snowflakeService.GenerateID()
	conversationOpportunityMembershipRequest.ConversationID = conversation.ID
	conversationOpportunityMembershipRequest.OpportunityMembershipRequestID = opportunityMembershipRequestID

	if err := s.conversationOpportunityMembershipRequestRepository.Create(conversationOpportunityMembershipRequest); err != nil {
		return 0, NewErrServerError()
	}

	// Message
	message := models.Message{}
	message.ConversationID = conversation.ID
	message.SenderID = volunteerID
	message.ID = s.snowflakeService.GenerateID()
	message.Type = models.MessageTypeVolunteerRequestProfile
	message.Timestamp = time.Now()
	message.Edited = false
	messageBody := MessageVolunteerRequestProfile{
		Message: messageStr,
		UserID:  volunteerID,
	}

	jsonBytes, err := json.Marshal(messageBody)
	if err != nil {
		return 0, NewErrServerError()
	}

	message.Body = postgres.Jsonb{
		RawMessage: json.RawMessage(jsonBytes),
	}

	if err := s.sendMessage(ctx, message); err != nil {
		s.logger.Error().Err(err).Msg("Error creating message")
		return 0, NewErrServerError()
	}

	return conversation.ID, nil
}

// GetUserConversationMemberships gets a user's volunteer conversation memberships.
func (s *service) GetUserConversationMemberships(userID int64) ([]models.ConversationMembership, error) {
	memberships, err := s.conversationMembershipRepository.FindByUserID(userID)
	if err != nil {
		return nil, NewErrServerError()
	}

	return memberships, nil
}

// ConversationsResponse contains conversations and total number of pages.
type ConversationsResponse struct {
	Conversations []ConversationView `json:"conversations"`
	Pages         uint               `json:"pages"`
}

// GetUserConversations gets all of a user's conversations.
func (s *service) GetUserConversations(ctx context.Context, userID int64) (*ConversationsResponse, error) {
	memberships, err := s.GetUserConversationMemberships(userID)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error getting user conversation memberships")
		return nil, NewErrServerError()
	}

	ids := []int64{}
	for _, membership := range memberships {
		ids = append(ids, membership.ConversationID)
	}

	res, err := s.conversationRepository.FindByIDs(ctx, ids)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error getting conversations by IDs")
		return nil, NewErrServerError()
	}

	views := []ConversationView{}

	for _, conversation := range res.Conversations {
		view := ConversationView{
			Conversation: conversation,
		}

		view.Conversation.Name = conversation.Organization.Name
		view.Conversation.ProfilePicture = conversation.Organization.ProfilePicture
		// Temporary dummy value
		view.UnreadCount = 0
		view.LastMessageView, _ = s.messageToView(ctx, conversation.LastMessage)

		views = append(views, view)
	}

	return &ConversationsResponse{
		Conversations: views,
		Pages:         uint(res.TotalResults/dbctx.Get(ctx).Limit) + 1,
	}, nil
}

// GetUserConversation gets a single conversation from a user perspective.
func (s *service) GetUserConversation(ctx context.Context, conversationID int64) (*ConversationView, error) {
	conversation, err := s.conversationRepository.FindByID(conversationID)
	if err != nil {
		return nil, NewErrConversationNotFound()
	}

	view := &ConversationView{
		Conversation: *conversation,
	}

	view.Conversation.Name = conversation.Organization.Name
	view.Conversation.ProfilePicture = conversation.Organization.ProfilePicture
	// Temporary dummy value
	view.UnreadCount = 0
	view.LastMessageView, _ = s.messageToView(ctx, conversation.LastMessage)

	requests, err := s.conversationOpportunityMembershipRequestRepository.FindByConversationID(conversationID)
	if err != nil {
		return nil, NewErrServerError()
	}

	for _, request := range requests {
		if request.OpportunityMembershipRequest == nil {
			s.logger.Warn().Msgf("ConversationOpportunityMembershipRequest %d linked to missing OpportunityMembershipRequest", request.ID)
			continue
		}

		view.Conversation.OpportunityMembershipRequests = append(conversation.OpportunityMembershipRequests, *request.OpportunityMembershipRequest)
	}

	return view, nil
}

// GetOrganizationConversations gets an organization's internal conversations.
func (s *service) GetOrganizationConversations(ctx context.Context, organizationID int64) (*ConversationsResponse, error) {
	res, err := s.conversationRepository.FindByOrganizationID(ctx, organizationID)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error getting conversations by organization ID")
		return nil, NewErrServerError()
	}

	views := []ConversationView{}

	for _, conversation := range res.Conversations {
		view := ConversationView{
			Conversation: conversation,
		}

		memberships, err := s.conversationMembershipRepository.FindByConversationID(conversation.ID)
		if err != nil {
			continue
		}

		if len(memberships) < 1 {
			continue
		}

		view.Conversation.Name = fmt.Sprintf("%s %s", memberships[0].User.FirstName, memberships[0].User.LastName)
		view.Conversation.ProfilePicture = memberships[0].User.ProfilePicture

		// Temporary dummy value
		view.UnreadCount = 0
		view.LastMessageView, _ = s.messageToView(ctx, conversation.LastMessage)

		views = append(views, view)
	}

	return &ConversationsResponse{
		Conversations: views,
		Pages:         uint(res.TotalResults/dbctx.Get(ctx).Limit) + 1,
	}, nil
}

// GetOrganizationConversation gets a single conversation from a conversation perspective.
func (s *service) GetOrganizationConversation(ctx context.Context, conversationID int64) (*ConversationView, error) {
	conversation, err := s.conversationRepository.FindByID(conversationID)
	if err != nil {
		return nil, NewErrConversationNotFound()
	}

	view := &ConversationView{
		Conversation: *conversation,
	}

	memberships, err := s.conversationMembershipRepository.FindByConversationID(conversation.ID)
	if err == nil && len(memberships) < 1 {
		view.Conversation.Name = fmt.Sprintf("%s %s", memberships[0].User.FirstName, memberships[0].User.LastName)
		view.Conversation.ProfilePicture = memberships[0].User.ProfilePicture
	}

	requests, err := s.conversationOpportunityMembershipRequestRepository.FindByConversationID(conversationID)
	if err != nil {
		return nil, NewErrServerError()
	}

	for _, request := range requests {
		if request.OpportunityMembershipRequest == nil {
			s.logger.Warn().Msgf("ConversationOpportunityMembershipRequest %d linked to missing OpportunityMembershipRequest", request.ID)
			continue
		}

		view.Conversation.OpportunityMembershipRequests = append(conversation.OpportunityMembershipRequests, *request.OpportunityMembershipRequest)
	}

	return view, nil
}

// SendStandardMessage sends a standard message to a conversation, returning the ID on success.
func (s *service) SendStandardMessage(ctx context.Context, conversationID, senderID int64, messageText string, asOrganization bool) (int64, error) {
	message := models.Message{}
	message.ID = s.snowflakeService.GenerateID()
	message.Timestamp = time.Now()
	message.ConversationID = conversationID
	message.SenderID = senderID
	message.Type = models.MessageTypeStandard
	perspective := models.MessageSenderPerspectiveVolunteer
	if asOrganization {
		perspective = models.MessageSenderPerspectiveOrganization
	}

	message.SenderPerspective = &perspective

	messageBody := MessageStandard{
		Text: messageText,
	}

	jsonBytes, err := marshalMessageBody(messageBody)
	if err != nil {
		return 0, NewErrServerError()
	}

	message.Body = *jsonBytes
	message.Edited = false
	if err := s.sendMessage(ctx, message); err != nil {
		s.logger.Error().Err(err).Msg("Error creating message")
		return 0, NewErrServerError()
	}

	return message.ID, nil
}

func (s *service) sendMessage(ctx context.Context, message models.Message) error {
	if err := s.messageRepository.Create(ctx, message); err != nil {
		s.logger.Error().Err(err).Msg("Error creating message")
		return err
	}

	go s.brokerPublishMessageSent(message)

	return nil
}

// brokerPublishMessageSent publishes a message as a MessageSent event.
// Should be called asynchronously/spawned as a goroutine.
func (s *service) brokerPublishMessageSent(message models.Message) error {
	view, err := s.messageToView(context.Background(), message)
	if err != nil {
		return err
	}

	if err := s.broker.Publish(stream, pubsub.Event{
		EventName: MessageSent,
		Payload:   *view,
	}); err != nil {
		s.logger.Error().Err(err).Msg("Error publishing message to pub/sub")
		return err
	}

	return nil
}

// ConversationMessagesResponse represents a response containing messages and paging information.
type ConversationMessagesResponse struct {
	Messages []MessageView `json:"messages"`
	Pages    uint          `json:"pages"`
}

// GetConversationMessages gets messages by conversation ID.
func (s *service) GetConversationMessages(ctx context.Context, conversationID int64) (*ConversationMessagesResponse, error) {
	res, err := s.messageRepository.FindByConversationID(ctx, conversationID)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error getting conversation messages")
		return nil, NewErrServerError()
	}

	views := []MessageView{}

	for _, message := range res.Messages {
		view, err := s.messageToView(ctx, message)
		if err != nil {
			s.logger.Error().Err(err).Msg("Error converting message to view")
			continue
		}

		views = append(views, *view)
	}

	return &ConversationMessagesResponse{
		Messages: views,
		Pages:    uint(res.TotalResults/dbctx.Get(ctx).Limit) + 1,
	}, nil
}

// messageToView converts a raw models.Message object into a *MessageView
// with a parsed body.
func (s *service) messageToView(ctx context.Context, message models.Message) (*MessageView, error) {
	body, err := s.parseMessage(ctx, message.Type, message.Body.RawMessage)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error parsing message body")
		return nil, err
	}

	view := &MessageView{
		ID:              message.ID,
		ConversationID:  message.ConversationID,
		SenderID:        message.SenderID,
		Timestamp:       message.Timestamp,
		Type:            message.Type,
		Edited:          message.Edited,
		EditedTimestamp: message.EditedTimestamp,
		Body:            body,
	}

	return view, nil
}

// parseMessage takes in a raw JSON message and adds necessary data to a message before returning it as an interface{}.
func (s *service) parseMessage(ctx context.Context, messageType string, rawMessage json.RawMessage) (interface{}, error) {
	switch messageType {
	case models.MessageTypeVolunteerRequestProfile:
		body := MessageVolunteerRequestProfile{}
		err := json.Unmarshal(rawMessage, &body)
		if err != nil {
			return nil, err
		}

		view, err := s.getMessageVolunteerRequestProfileView(ctx, body.UserID)
		if err != nil {
			return nil, err
		}

		return view, nil
	case models.MessageTypeVolunteerRequestAcceptance:
		body := MessageTypeVolunteerRequestAcceptance{}
		err := json.Unmarshal(rawMessage, &body)
		if err != nil {
			return nil, err
		}

		view, err := s.getMessageVolunteerRequestAcceptance(ctx, body.UserID, body.OpportunityID)
		if err != nil {
			return nil, err
		}

		return view, nil
	case models.MessageTypeHoursRequested:
		body := MessageTypeHoursRequested{}
		err := json.Unmarshal(rawMessage, &body)
		if err != nil {
			return nil, err
		}

		view, err := s.getMessageTypeHoursRequestedView(ctx, body.VolunteeringHourLogRequestID)
		if err != nil {
			return nil, err
		}

		return view, nil
	case models.MessageTypeHoursAccepted:
		body := MessageTypeHoursAccepted{}
		err := json.Unmarshal(rawMessage, &body)
		if err != nil {
			return nil, err
		}

		view, err := s.getMessageTypeHoursRequestedView(ctx, body.VolunteeringHourLogRequestID)
		if err != nil {
			return nil, err
		}

		return view, nil
	case models.MessageTypeHoursDeclined:
		body := MessageTypeHoursDeclined{}
		err := json.Unmarshal(rawMessage, &body)
		if err != nil {
			return nil, err
		}

		view, err := s.getMessageTypeHoursRequestedView(ctx, body.VolunteeringHourLogRequestID)
		if err != nil {
			return nil, err
		}

		return view, nil
	}

	// Fallback
	body := map[string]interface{}{}
	err := json.Unmarshal(rawMessage, &body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// getMessageVolunteerRequestProfileView gets a MessageVolunteerRequestProfileView for a user by ID.
func (s *service) getMessageVolunteerRequestProfileView(ctx context.Context, userID int64) (*MessageVolunteerRequestProfileView, error) {
	profile := &MessageVolunteerRequestProfileView{}
	// Find the user to verify that it is active.
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, NewErrUserNotFound()
	}

	profile.UserID = user.ID
	profile.FirstName = user.FirstName
	profile.LastName = user.LastName
	profile.ProfilePicture = user.ProfilePicture
	profile.DateOfBirth = user.DateOfBirth

	// Find all UserTag objects by UserID.
	userTags, err := s.userTagRepository.FindByUserID(userID)
	if err != nil {
		return nil, NewErrServerError()
	}

	tags := []models.Tag{}
	for _, userTag := range userTags {
		// Get the tag by ID.
		tag, err := s.tagRepository.FindByID(userTag.TagID)
		if err != nil {
			// Tag not found, skip.
			s.logger.Error().Err(err).Msg("Error getting user tags: UserTag object missing valid Tag")
			continue
		}

		// Append the tag to the tags array.
		tags = append(tags, *tag)
	}

	profile.Tags = tags

	// Location
	if user.LocationLatitude != 0.0 || user.LocationLongitude != 0.0 {
		coordinates := &location.Coordinates{
			Latitude:  user.LocationLatitude,
			Longitude: user.LocationLongitude,
		}

		location, err := s.locationService.CoordinatesToCity(coordinates)
		if err == nil {
			profile.Location = location
		}
	}

	// Profile fields
	fields, err := s.userProfileFieldRepository.FindByUserID(userID)
	if err != nil {
		return nil, NewErrServerError()
	}

	profile.ProfileFields = fields

	profile.PreviousExperience = &PreviousExperience{
		Count: 0,
	}

	return profile, nil
}

// getMessageVolunteerRequestAcceptance gets a MessageTypeVolunteerRequestAcceptanceView by user ID and opportunity ID.
func (s *service) getMessageVolunteerRequestAcceptance(ctx context.Context, userID, opportunityID int64) (*MessageTypeVolunteerRequestAcceptanceView, error) {
	view := &MessageTypeVolunteerRequestAcceptanceView{}

	opportunity, err := s.opportunityRepository.FindByID(ctx, opportunityID)
	if err != nil {
		return nil, err
	}

	view.UserID = userID
	view.OpportunityID = opportunityID
	view.OpportunityTitle = opportunity.Title

	return view, nil
}

// getMessageTypeHoursRequestedView gets a MessageTypeHoursRequestedView by request ID.
func (s *service) getMessageTypeHoursRequestedView(ctx context.Context, requestID int64) (*MessageTypeHoursRequestedView, error) {
	view := &MessageTypeHoursRequestedView{}

	request, err := s.volunteeringHourLogRequestRepository.FindByID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	view.VolunteeringHourLogRequest = *request

	return view, nil
}

// marshalMessageBody marshals an interface of a message body to a postgres Jsonb value.
func marshalMessageBody(body interface{}) (*postgres.Jsonb, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return &postgres.Jsonb{
		RawMessage: json.RawMessage(jsonBytes),
	}, nil
}

// SendHoursRequestMessage sends an hours request message to a user's organization message.
func (s *service) SendHoursRequestMessage(ctx context.Context, userID, organizationID, requestID int64) (int64, error) {
	conversation, err := s.conversationRepository.FindUserOrganizationConversation(ctx, userID, organizationID)
	if err != nil {
		return 0, NewErrConversationNotFound()
	}

	message := models.Message{}
	message.ID = s.snowflakeService.GenerateID()
	message.Timestamp = time.Now()
	message.ConversationID = conversation.ID
	message.SenderID = userID
	message.Type = models.MessageTypeHoursRequested
	perspective := models.MessageSenderPerspectiveVolunteer

	message.SenderPerspective = &perspective

	messageBody := MessageTypeHoursRequested{
		VolunteeringHourLogRequestID: requestID,
	}

	jsonBytes, err := marshalMessageBody(messageBody)
	if err != nil {
		return 0, NewErrServerError()
	}

	message.Body = *jsonBytes
	message.Edited = false
	if err := s.sendMessage(ctx, message); err != nil {
		s.logger.Error().Err(err).Msg("Error creating message")
		return 0, NewErrServerError()
	}

	return message.ID, nil
}

// SendHoursRequestAcceptedMessage sends an hours request accept message to a user's organization message.
func (s *service) SendHoursRequestAcceptedMessage(ctx context.Context, userID, requestID int64) (int64, error) {
	request, err := s.volunteeringHourLogRequestRepository.FindByID(ctx, requestID)
	if err != nil {
		return 0, NewErrServerError()
	}

	conversation, err := s.conversationRepository.FindUserOrganizationConversation(ctx, request.VolunteerID, request.OrganizationID)
	if err != nil {
		return 0, NewErrConversationNotFound()
	}

	message := models.Message{}
	message.ID = s.snowflakeService.GenerateID()
	message.Timestamp = time.Now()
	message.ConversationID = conversation.ID
	message.SenderID = userID
	message.Type = models.MessageTypeHoursAccepted
	perspective := models.MessageSenderPerspectiveVolunteer

	message.SenderPerspective = &perspective

	messageBody := MessageTypeHoursAccepted{
		VolunteeringHourLogRequestID: requestID,
	}

	jsonBytes, err := marshalMessageBody(messageBody)
	if err != nil {
		return 0, NewErrServerError()
	}

	message.Body = *jsonBytes
	message.Edited = false
	if err := s.sendMessage(ctx, message); err != nil {
		s.logger.Error().Err(err).Msg("Error creating message")
		return 0, NewErrServerError()
	}

	return message.ID, nil
}

// SendHoursRequestDeclinedMessage sends an hours request accept message to a user's organization message.
func (s *service) SendHoursRequestDeclinedMessage(ctx context.Context, userID, requestID int64) (int64, error) {
	request, err := s.volunteeringHourLogRequestRepository.FindByID(ctx, requestID)
	if err != nil {
		return 0, NewErrServerError()
	}

	conversation, err := s.conversationRepository.FindUserOrganizationConversation(ctx, request.VolunteerID, request.OrganizationID)
	if err != nil {
		return 0, NewErrConversationNotFound()
	}

	message := models.Message{}
	message.ID = s.snowflakeService.GenerateID()
	message.Timestamp = time.Now()
	message.ConversationID = conversation.ID
	message.SenderID = userID
	message.Type = models.MessageTypeHoursDeclined
	perspective := models.MessageSenderPerspectiveVolunteer

	message.SenderPerspective = &perspective

	messageBody := MessageTypeHoursDeclined{
		VolunteeringHourLogRequestID: requestID,
	}

	jsonBytes, err := marshalMessageBody(messageBody)
	if err != nil {
		return 0, NewErrServerError()
	}

	message.Body = *jsonBytes
	message.Edited = false
	if err := s.sendMessage(ctx, message); err != nil {
		s.logger.Error().Err(err).Msg("Error creating message")
		return 0, NewErrServerError()
	}

	return message.ID, nil
}

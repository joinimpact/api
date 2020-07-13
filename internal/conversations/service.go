package conversations

import (
	"context"
	"encoding/json"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joinimpact/api/internal/cdn"
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/email"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/snowflakes"
	"github.com/joinimpact/api/internal/users"
	"github.com/rs/zerolog"
)

// Service defines methods for interacting with conversations and messages.
type Service interface {
	// CreateOpportunityMembershipRequestConversation creates an opportunity membership request conversation and adds a message to it. Returns conversation ID on success.
	CreateOpportunityMembershipRequestConversation(ctx context.Context, organizationID, opportunityID, opportunityMembershipRequestID, volunteerID int64, messageStr string) (int64, error)
}

// service represents the internal implementation of the conversations Service.
type service struct {
	conversationRepository                             models.ConversationRepository
	conversationMembershipRepository                   models.ConversationMembershipRepository
	conversationOpportunityMembershipRequestRepository models.ConversationOpportunityMembershipRequestRepository
	conversationOrganizationMembershipRepository       models.ConversationOrganizationMembershipRepository
	messageRepository                                  models.MessageRepository
	usersService                                       users.Service // TODO: find a way to use the users service without a dependency like this
	config                                             *config.Config
	logger                                             *zerolog.Logger
	snowflakeService                                   snowflakes.SnowflakeService
	emailService                                       email.Service
	cdnClient                                          *cdn.Client
}

// NewService creates and returns a new conversations.Service.
func NewService(conversationRepository models.ConversationRepository, conversationMembershipRepository models.ConversationMembershipRepository, conversationOpportunityMembershipRequestRepository models.ConversationOpportunityMembershipRequestRepository, conversationOrganizationMembershipRepository models.ConversationOrganizationMembershipRepository, messageRepository models.MessageRepository, usersService users.Service, config *config.Config, logger *zerolog.Logger, snowflakeService snowflakes.SnowflakeService, emailService email.Service) Service {
	return &service{
		conversationRepository,
		conversationMembershipRepository,
		conversationOpportunityMembershipRequestRepository,
		conversationOrganizationMembershipRepository,
		messageRepository,
		usersService,
		config,
		logger,
		snowflakeService,
		emailService,
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

	if err := s.messageRepository.Create(ctx, message); err != nil {
		return 0, NewErrServerError()
	}

	return conversation.ID, nil
}

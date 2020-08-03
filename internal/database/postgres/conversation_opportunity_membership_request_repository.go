package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// conversationOpportunityMembershipRequestRepository stores and controls conversationOpportunityMembershipRequests in the database.
type conversationOpportunityMembershipRequestRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewConversationOpportunityMembershipRequestRepository creates and returns a new ConversationOpportunityMembershipRequestRepository.
func NewConversationOpportunityMembershipRequestRepository(db *gorm.DB, logger *zerolog.Logger) models.ConversationOpportunityMembershipRequestRepository {
	return &conversationOpportunityMembershipRequestRepository{db, logger}
}

// FindByID finds a single entity by ID.
func (r *conversationOpportunityMembershipRequestRepository) FindByID(id int64) (*models.ConversationOpportunityMembershipRequest, error) {
	var conversationOpportunityMembershipRequest models.ConversationOpportunityMembershipRequest
	if err := r.db.First(&conversationOpportunityMembershipRequest, id).Error; err != nil {
		return &conversationOpportunityMembershipRequest, err
	}
	return &conversationOpportunityMembershipRequest, nil
}

// FindByConversationID finds multiple entities by the conversation ID.
func (r *conversationOpportunityMembershipRequestRepository) FindByConversationID(conversationID int64) ([]models.ConversationOpportunityMembershipRequest, error) {
	var opportunities []models.ConversationOpportunityMembershipRequest
	if err := r.db.Preload("OpportunityMembershipRequest").Where("conversation_id = ?", conversationID).Find(&opportunities).Error; err != nil {
		return opportunities, err
	}
	return opportunities, nil
}

// Create creates a new User.
func (r *conversationOpportunityMembershipRequestRepository) Create(conversationOpportunityMembershipRequest models.ConversationOpportunityMembershipRequest) error {
	return r.db.Create(&conversationOpportunityMembershipRequest).Error
}

// Update updates a User with the ID in the provided User.
func (r *conversationOpportunityMembershipRequestRepository) Update(conversationOpportunityMembershipRequest models.ConversationOpportunityMembershipRequest) error {
	return r.db.Model(&models.ConversationOpportunityMembershipRequest{}).Updates(conversationOpportunityMembershipRequest).Error
}

// DeleteByID deletes a User by ID.
func (r *conversationOpportunityMembershipRequestRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.ConversationOpportunityMembershipRequest{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

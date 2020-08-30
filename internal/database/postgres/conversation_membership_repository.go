package postgres

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// conversationMembershipRepository stores and controls conversationMemberships in the database.
type conversationMembershipRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewConversationMembershipRepository creates and returns a new ConversationMembershipRepository.
func NewConversationMembershipRepository(db *gorm.DB, logger *zerolog.Logger) models.ConversationMembershipRepository {
	return &conversationMembershipRepository{db, logger}
}

// FindByID finds a single entity by ID.
func (r *conversationMembershipRepository) FindByID(id int64) (*models.ConversationMembership, error) {
	var conversationMembership models.ConversationMembership
	if err := r.db.First(&conversationMembership, id).Error; err != nil {
		return &conversationMembership, err
	}
	return &conversationMembership, nil
}

// FindByConversationID finds multiple entities by the conversation ID.
func (r *conversationMembershipRepository) FindByConversationID(conversationID int64) ([]models.ConversationMembership, error) {
	var conversationMemberships []models.ConversationMembership
	if err := r.db.Preload("User").Where("conversation_id = ? AND active = True", conversationID).Find(&conversationMemberships).Error; err != nil {
		return conversationMemberships, err
	}
	return conversationMemberships, nil
}

// FindByUserID finds multiple entities by the creator ID.
func (r *conversationMembershipRepository) FindByUserID(userID int64) ([]models.ConversationMembership, error) {
	var conversationMemberships []models.ConversationMembership
	if err := r.db.Where("user_id = ? AND active = True", userID).Find(&conversationMemberships).Error; err != nil {
		return conversationMemberships, err
	}
	return conversationMemberships, nil
}

// FindByUserIDAndConversationID finds a single entity by user ID and conversation ID.
func (r *conversationMembershipRepository) FindByUserIDAndConversationID(ctx context.Context, userID, conversationID int64) (*models.ConversationMembership, error) {
	var conversationMembership models.ConversationMembership
	if err := r.db.Where("user_id = ? AND conversation_id = ? AND active = True", userID, conversationID).First(&conversationMembership).Error; err != nil {
		return &conversationMembership, err
	}
	return &conversationMembership, nil

}

// Create creates a new User.
func (r *conversationMembershipRepository) Create(conversationMembership models.ConversationMembership) error {
	return r.db.Create(&conversationMembership).Error
}

// Update updates a User with the ID in the provided User.
func (r *conversationMembershipRepository) Update(conversationMembership models.ConversationMembership) error {
	return r.db.Model(&models.ConversationMembership{}).Updates(conversationMembership).Error
}

// DeleteByID deletes a User by ID.
func (r *conversationMembershipRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.ConversationMembership{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

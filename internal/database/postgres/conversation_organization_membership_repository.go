package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// conversationOrganizationMembershipRepository stores and controls conversationOrganizationMemberships in the database.
type conversationOrganizationMembershipRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewConversationOrganizationMembershipRepository creates and returns a new ConversationOrganizationMembershipRepository.
func NewConversationOrganizationMembershipRepository(db *gorm.DB, logger *zerolog.Logger) models.ConversationOrganizationMembershipRepository {
	return &conversationOrganizationMembershipRepository{db, logger}
}

// FindByID finds a single entity by ID.
func (r *conversationOrganizationMembershipRepository) FindByID(id int64) (*models.ConversationOrganizationMembership, error) {
	var conversationOrganizationMembership models.ConversationOrganizationMembership
	if err := r.db.First(&conversationOrganizationMembership, id).Error; err != nil {
		return &conversationOrganizationMembership, err
	}
	return &conversationOrganizationMembership, nil
}

// FindByConversationID finds multiple entities by the conversation ID.
func (r *conversationOrganizationMembershipRepository) FindByConversationID(conversationID int64) ([]models.ConversationOrganizationMembership, error) {
	var conversationOrganizationMemberships []models.ConversationOrganizationMembership
	if err := r.db.Where("conversation_id = ? AND active = True", conversationID).Find(&conversationOrganizationMemberships).Error; err != nil {
		return conversationOrganizationMemberships, err
	}
	return conversationOrganizationMemberships, nil
}

// FindByOrganizationID finds multiple entities by the creator ID.
func (r *conversationOrganizationMembershipRepository) FindByOrganizationID(organizationID int64) ([]models.ConversationOrganizationMembership, error) {
	var conversationOrganizationMemberships []models.ConversationOrganizationMembership
	if err := r.db.Where("organization_id = ? AND active = True", organizationID).Find(&conversationOrganizationMemberships).Error; err != nil {
		return conversationOrganizationMemberships, err
	}
	return conversationOrganizationMemberships, nil
}

// Create creates a new User.
func (r *conversationOrganizationMembershipRepository) Create(conversationOrganizationMembership models.ConversationOrganizationMembership) error {
	return r.db.Create(&conversationOrganizationMembership).Error
}

// Update updates a User with the ID in the provided User.
func (r *conversationOrganizationMembershipRepository) Update(conversationOrganizationMembership models.ConversationOrganizationMembership) error {
	return r.db.Model(&models.ConversationOrganizationMembership{}).Updates(conversationOrganizationMembership).Error
}

// DeleteByID deletes a User by ID.
func (r *conversationOrganizationMembershipRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.ConversationOrganizationMembership{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

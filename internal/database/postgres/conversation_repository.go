package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// conversationRepository stores and controls conversations in the database.
type conversationRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewConversationRepository creates and returns a new ConversationRepository.
func NewConversationRepository(db *gorm.DB, logger *zerolog.Logger) models.ConversationRepository {
	return &conversationRepository{db, logger}
}

// FindByID finds a single entity by ID.
func (r *conversationRepository) FindByID(id int64) (*models.Conversation, error) {
	var conversation models.Conversation
	if err := r.db.First(&conversation, id).Error; err != nil {
		return &conversation, err
	}
	return &conversation, nil
}

// FindByOrganizationID finds multiple entities by the organization ID.
func (r *conversationRepository) FindByOrganizationID(organizationID int64) ([]models.Conversation, error) {
	var opportunities []models.Conversation
	if err := r.db.Where("organization_id = ? AND active = True", organizationID).Find(&opportunities).Error; err != nil {
		return opportunities, err
	}
	return opportunities, nil
}

// FindByCreatorID finds multiple entities by the creator ID.
func (r *conversationRepository) FindByCreatorID(creatorID int64) ([]models.Conversation, error) {
	var opportunities []models.Conversation
	if err := r.db.Where("creator_id = ? AND active = True", creatorID).Find(&opportunities).Error; err != nil {
		return opportunities, err
	}
	return opportunities, nil
}

// Create creates a new User.
func (r *conversationRepository) Create(conversation models.Conversation) error {
	return r.db.Create(&conversation).Error
}

// Update updates a User with the ID in the provided User.
func (r *conversationRepository) Update(conversation models.Conversation) error {
	return r.db.Model(&models.Conversation{}).Updates(conversation).Error
}

// DeleteByID deletes a User by ID.
func (r *conversationRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.Conversation{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

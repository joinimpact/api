package postgres

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/pkg/dbctx"
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
	if err := r.db.Preload("LastMessage", func(db *gorm.DB) *gorm.DB {
		return db.Order("timestamp desc")
	}).Preload("Organization").First(&conversation, id).Error; err != nil {
		return &conversation, err
	}
	return &conversation, nil
}

// FindByIDs finds multiple entities by IDs.
func (r *conversationRepository) FindByIDs(ctx context.Context, ids []int64) (*models.ConversationsResponse, error) {
	response := &models.ConversationsResponse{}

	dbctx := dbctx.Get(ctx)

	db := r.db.
		Model(&models.Conversation{}).
		Preload("LastMessage", func(db *gorm.DB) *gorm.DB {
			return db.Order("timestamp desc")
		}).
		Preload("Organization").
		Limit(dbctx.Limit).
		Joins("LEFT JOIN (select distinct on (timestamp) * from messages order by timestamp desc limit 1) as message ON message.conversation_id = conversations.id").
		Where("conversations.id IN (?) AND active = True", ids).
		Order("message.timestamp asc").
		Count(&response.TotalResults).
		Offset(dbctx.Page * dbctx.Limit)

	if err := db.Find(&response.Conversations).Error; err != nil {
		return response, err
	}

	return response, nil
}

// FindByOrganizationID finds multiple entities by the organization ID.
func (r *conversationRepository) FindByOrganizationID(ctx context.Context, organizationID int64) (*models.ConversationsResponse, error) {
	response := &models.ConversationsResponse{}

	dbctx := dbctx.Get(ctx)

	db := r.db.
		Model(&models.Conversation{}).
		Preload("LastMessage", func(db *gorm.DB) *gorm.DB {
			return db.Order("timestamp desc")
		}).
		Limit(dbctx.Limit).
		Joins("LEFT JOIN (select distinct on (timestamp) * from messages order by timestamp desc limit 1) as message ON message.conversation_id = conversations.id").
		Where("conversations.organization_id = ? AND active = True", organizationID).
		Order("message.timestamp asc").
		Count(&response.TotalResults).
		Offset(dbctx.Page * dbctx.Limit)

	if err := db.Find(&response.Conversations).Error; err != nil {
		return response, err
	}

	return response, nil
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

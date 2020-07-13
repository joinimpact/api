package postgres

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// messageRepository stores and controls messages in the database.
type messageRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewMessageRepository creates and returns a new MessageRepository.
func NewMessageRepository(db *gorm.DB, logger *zerolog.Logger) models.MessageRepository {
	return &messageRepository{db, logger}
}

// FindByID finds a single entity by ID.
func (r *messageRepository) FindByID(ctx context.Context, id int64) (*models.Message, error) {
	var message models.Message
	if err := r.db.First(&message, id).Error; err != nil {
		return &message, err
	}
	return &message, nil
}

// FindByConversationID finds multiple entities by the conversation ID.
func (r *messageRepository) FindByConversationID(ctx context.Context, conversationID int64) ([]models.Message, error) {
	var messages []models.Message
	if err := r.db.Where("conversation_id = ?", conversationID).Order("timestamp DESC", true).Find(&messages).Error; err != nil {
		return messages, err
	}
	return messages, nil
}

// FindBySenderID finds multiple entities by the sender ID.
func (r *messageRepository) FindBySenderID(ctx context.Context, senderID int64) ([]models.Message, error) {
	var messages []models.Message
	if err := r.db.Where("sender_id = ?", senderID).Order("timestamp DESC", true).Find(&messages).Error; err != nil {
		return messages, err
	}
	return messages, nil
}

// FindInConversationBySenderID finds multiple entities by the conversation and sender ID.
func (r *messageRepository) FindInConversationBySenderID(ctx context.Context, conversationID, senderID int64) ([]models.Message, error) {
	var messages []models.Message
	if err := r.db.Where("conversation_id = ? AND sender_id = ?", conversationID, senderID).Order("timestamp DESC", true).Find(&messages).Error; err != nil {
		return messages, err
	}
	return messages, nil
}

// Create creates a new User.
func (r *messageRepository) Create(ctx context.Context, message models.Message) error {
	return r.db.Create(&message).Error
}

// Update updates a User with the ID in the provided User.
func (r *messageRepository) Update(ctx context.Context, message models.Message) error {
	return r.db.Model(&models.Message{}).Updates(message).Error
}

// DeleteByID deletes a User by ID.
func (r *messageRepository) DeleteByID(ctx context.Context, id int64) error {
	return r.db.Delete(&models.Message{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

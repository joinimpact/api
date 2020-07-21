package postgres

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/pkg/dbctx"
	"github.com/rs/zerolog"
)

// eventResponseRepository stores and controls events in the database.
type eventResponseRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewEventResponseRepository creates and returns a new EventResponseRepository.
func NewEventResponseRepository(db *gorm.DB, logger *zerolog.Logger) models.EventResponseRepository {
	return &eventResponseRepository{db, logger}
}

// FindByID finds a single entity by ID.
func (r *eventResponseRepository) FindByID(ctx context.Context, id int64) (*models.EventResponse, error) {
	var eventResponse models.EventResponse
	if err := r.db.First(&eventResponse, id).Error; err != nil {
		return &eventResponse, err
	}
	return &eventResponse, nil
}

// FindByOpportunityID finds multiple entities by the opportunity ID.
func (r *eventResponseRepository) FindByOpportunityID(ctx context.Context, opportunityID int64) ([]models.EventResponse, error) {
	var eventResponses []models.EventResponse
	if err := r.db.
		Limit(dbctx.Get(ctx).Limit).
		Offset(dbctx.Get(ctx).Page*dbctx.Get(ctx).Limit).
		Where("opportunity_id = ?", opportunityID).
		Find(&eventResponses).
		Error; err != nil {
		return eventResponses, err
	}
	return eventResponses, nil
}

// FindByEventID finds multiple entities by the event ID.
func (r *eventResponseRepository) FindByEventID(ctx context.Context, eventID int64) ([]models.EventResponse, error) {
	var eventResponses []models.EventResponse
	if err := r.db.Where("event_id = ?", eventID).Find(&eventResponses).Error; err != nil {
		return eventResponses, err
	}
	return eventResponses, nil
}

// FindByUserID finds multiple entities by the user ID.
func (r *eventResponseRepository) FindByUserID(ctx context.Context, userID int64) ([]models.EventResponse, error) {
	var eventResponses []models.EventResponse
	if err := r.db.Where("user_id = ?", userID).Find(&eventResponses).Error; err != nil {
		return eventResponses, err
	}
	return eventResponses, nil
}

// FindInEventByUserID finds an entity by the event and user ID.
func (r *eventResponseRepository) FindInEventByUserID(ctx context.Context, eventID, userID int64) (*models.EventResponse, error) {
	var eventResponse models.EventResponse
	if err := r.db.Where("event_id = ? AND user_id = ?", eventID, userID).First(&eventResponse).Error; err != nil {
		return &eventResponse, err
	}
	return &eventResponse, nil
}

// Create creates a new Event.
func (r *eventResponseRepository) Create(ctx context.Context, event models.EventResponse) error {
	return r.db.Create(&event).Error
}

// Update updates a Event with the ID in the provided Event.
func (r *eventResponseRepository) Update(ctx context.Context, event models.EventResponse) error {
	return r.db.Model(&models.EventResponse{}).Updates(event).Error
}

// Save saves all fields in the provided entity.
func (r *eventResponseRepository) Save(ctx context.Context, event models.EventResponse) error {
	return r.db.Save(event).Error
}

// DeleteByID deletes a Event by ID.
func (r *eventResponseRepository) DeleteByID(ctx context.Context, id int64) error {
	return r.db.Delete(&models.EventResponse{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

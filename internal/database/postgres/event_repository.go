package postgres

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/pkg/dbctx"
	"github.com/rs/zerolog"
)

// eventRepository stores and controls events in the database.
type eventRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewEventRepository creates and returns a new EventRepository.
func NewEventRepository(db *gorm.DB, logger *zerolog.Logger) models.EventRepository {
	return &eventRepository{db, logger}
}

// FindByID finds a single entity by ID.
func (r *eventRepository) FindByID(ctx context.Context, id int64) (*models.Event, error) {
	var event models.Event
	if err := r.db.First(&event, id).Error; err != nil {
		return &event, err
	}
	return &event, nil
}

// FindByOpportunityID finds multiple entities by the opportunity ID.
func (r *eventRepository) FindByOpportunityID(ctx context.Context, opportunityID int64) ([]models.Event, error) {
	var events []models.Event

	dbctx := dbctx.Get(ctx)

	db := r.db.
		Limit(dbctx.Limit).
		Offset(dbctx.Page*dbctx.Limit).
		Where("opportunity_id = ? AND active = True AND title LIKE ?", opportunityID, fmt.Sprintf("%%%s%%", dbctx.Query)).
		Order("from_date ASC", true)

	if dbctx.From != nil {
		db = db.Where("from_date >= ?", *dbctx.From)
	}

	if dbctx.To != nil {
		db = db.Where("to_date < ?", *dbctx.To)
	}

	if err := db.
		Find(&events).
		Error; err != nil {
		return events, err
	}
	return events, nil
}

// FindByOpportunityIDs finds multiple entities by multiple opportunity IDs.
func (r *eventRepository) FindByOpportunityIDs(ctx context.Context, opportunityIDs []int64) ([]models.Event, error) {
	var events []models.Event

	dbctx := dbctx.Get(ctx)

	db := r.db.
		Limit(dbctx.Limit).
		Offset(dbctx.Page*dbctx.Limit).
		Where("opportunity_id IN (?) AND active = True AND title LIKE ?", opportunityIDs, fmt.Sprintf("%%%s%%", dbctx.Query)).
		Order("from_date ASC", true)

	if dbctx.From != nil {
		db = db.Where("from_date >= ? OR to_date >= ?", *dbctx.From, *dbctx.From)
	}

	if dbctx.To != nil {
		db = db.Where("from_date < ? OR to_date < ?", *dbctx.To, *dbctx.To)
	}

	if err := db.
		Find(&events).
		Error; err != nil {
		return events, err
	}
	return events, nil
}

// FindByCreatorID finds multiple entities by the creator ID.
func (r *eventRepository) FindByCreatorID(ctx context.Context, creatorID int64) ([]models.Event, error) {
	var events []models.Event
	if err := r.db.Where("creator_id = ? AND active = True", creatorID).Find(&events).Error; err != nil {
		return events, err
	}
	return events, nil
}

// Create creates a new Event.
func (r *eventRepository) Create(ctx context.Context, event models.Event) error {
	return r.db.Create(&event).Error
}

// Update updates a Event with the ID in the provided Event.
func (r *eventRepository) Update(ctx context.Context, event models.Event) error {
	return r.db.Model(&models.Event{}).Updates(event).Error
}

// Save saves all fields in the provided entity.
func (r *eventRepository) Save(ctx context.Context, event models.Event) error {
	return r.db.Save(event).Error
}

// DeleteByID deletes a Event by ID.
func (r *eventRepository) DeleteByID(ctx context.Context, id int64) error {
	return r.db.Delete(&models.Event{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

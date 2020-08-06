package postgres

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/pkg/dbctx"
	"github.com/rs/zerolog"
)

// volunteeringHourLogRepository stores and controls conversations in the database.
type volunteeringHourLogRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewVolunteeringHourLogRepository creates and returns a new VolunteeringHourLogRepository.
func NewVolunteeringHourLogRepository(db *gorm.DB, logger *zerolog.Logger) models.VolunteeringHourLogRepository {
	return &volunteeringHourLogRepository{db, logger}
}

// FindByID finds a single entity by ID.
func (r *volunteeringHourLogRepository) FindByID(ctx context.Context, id int64) (*models.VolunteeringHourLog, error) {
	var volunteeringHourLog models.VolunteeringHourLog
	if err := r.db.Preload("Organization").First(&volunteeringHourLog, id).Error; err != nil {
		return &volunteeringHourLog, err
	}
	return &volunteeringHourLog, nil
}

// FindByIDs finds multiple entities by IDs.
func (r *volunteeringHourLogRepository) FindByIDs(ctx context.Context, ids []int64) (*models.VolunteeringHourLogsResponse, error) {
	response := &models.VolunteeringHourLogsResponse{}

	dbctx := dbctx.Get(ctx)

	db := r.db.
		Model(&models.VolunteeringHourLog{}).
		Limit(dbctx.Limit).
		Where("id IN (?) AND active = True", ids).
		Count(&response.TotalResults).
		Offset(dbctx.Page * dbctx.Limit)

	if err := db.Find(&response.VolunteeringHourLogs).Error; err != nil {
		return response, err
	}

	return response, nil
}

// FindByOrganizationID finds multiple entities by the organization ID.
func (r *volunteeringHourLogRepository) FindByOrganizationID(ctx context.Context, organizationID int64) (*models.VolunteeringHourLogsResponse, error) {
	response := &models.VolunteeringHourLogsResponse{}

	dbctx := dbctx.Get(ctx)

	db := r.db.
		Model(&models.VolunteeringHourLog{}).
		Limit(dbctx.Limit).
		Where("organization_id = ?", organizationID).
		Count(&response.TotalResults).
		Offset(dbctx.Page * dbctx.Limit)

	if err := db.Find(&response.VolunteeringHourLogs).Error; err != nil {
		return response, err
	}

	return response, nil
}

// FindByOpportunityID finds multiple entities by the opportunity ID.
func (r *volunteeringHourLogRepository) FindByOpportunityID(ctx context.Context, opportunityID int64) (*models.VolunteeringHourLogsResponse, error) {
	response := &models.VolunteeringHourLogsResponse{}

	dbctx := dbctx.Get(ctx)

	db := r.db.
		Model(&models.VolunteeringHourLog{}).
		Limit(dbctx.Limit).
		Where("opportunity_id = ?", opportunityID).
		Count(&response.TotalResults).
		Offset(dbctx.Page * dbctx.Limit)

	if err := db.Find(&response.VolunteeringHourLogs).Error; err != nil {
		return response, err
	}

	return response, nil
}

// FindByVolunteerID finds multiple entities by the volunteer ID.
func (r *volunteeringHourLogRepository) FindByVolunteerID(ctx context.Context, volunteerID int64) (*models.VolunteeringHourLogsResponse, error) {
	response := &models.VolunteeringHourLogsResponse{}

	dbctx := dbctx.Get(ctx)

	db := r.db.
		Model(&models.VolunteeringHourLog{}).
		Limit(dbctx.Limit).
		Where("volunteer_id = ?", volunteerID).
		Count(&response.TotalResults).
		Offset(dbctx.Page * dbctx.Limit)

	if err := db.Find(&response.VolunteeringHourLogs).Error; err != nil {
		return response, err
	}

	return response, nil
}

// Create creates a new User.
func (r *volunteeringHourLogRepository) Create(ctx context.Context, volunteeringHourLog models.VolunteeringHourLog) error {
	return r.db.Create(&volunteeringHourLog).Error
}

// Update updates a User with the ID in the provided User.
func (r *volunteeringHourLogRepository) Update(ctx context.Context, volunteeringHourLog models.VolunteeringHourLog) error {
	return r.db.Model(&models.VolunteeringHourLog{}).Updates(volunteeringHourLog).Error
}

// DeleteByID deletes a User by ID.
func (r *volunteeringHourLogRepository) DeleteByID(ctx context.Context, id int64) error {
	return r.db.Delete(&models.VolunteeringHourLog{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

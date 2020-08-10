package postgres

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/pkg/dbctx"
	"github.com/rs/zerolog"
)

// volunteeringHourLogRequestRepository stores and controls conversations in the database.
type volunteeringHourLogRequestRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewVolunteeringHourLogRequestRepository creates and returns a new VolunteeringHourLogRequestRepository.
func NewVolunteeringHourLogRequestRepository(db *gorm.DB, logger *zerolog.Logger) models.VolunteeringHourLogRequestRepository {
	return &volunteeringHourLogRequestRepository{db, logger}
}

// FindByID finds a single entity by ID.
func (r *volunteeringHourLogRequestRepository) FindByID(ctx context.Context, id int64) (*models.VolunteeringHourLogRequest, error) {
	var volunteeringHourLogRequest models.VolunteeringHourLogRequest
	if err := r.db.Preload("Organization").First(&volunteeringHourLogRequest, id).Error; err != nil {
		return &volunteeringHourLogRequest, err
	}
	return &volunteeringHourLogRequest, nil
}

// FindByIDs finds multiple entities by IDs.
func (r *volunteeringHourLogRequestRepository) FindByIDs(ctx context.Context, ids []int64) (*models.VolunteeringHourLogRequestsResponse, error) {
	response := &models.VolunteeringHourLogRequestsResponse{}

	dbctx := dbctx.Get(ctx)

	db := r.db.
		Model(&models.VolunteeringHourLogRequest{}).
		Limit(dbctx.Limit).
		Where("id IN (?) AND active = True", ids).
		Count(&response.TotalResults).
		Offset(dbctx.Page * dbctx.Limit)

	if err := db.Find(&response.VolunteeringHourLogRequests).Error; err != nil {
		return response, err
	}

	return response, nil
}

// FindByOrganizationID finds multiple entities by the organization ID.
func (r *volunteeringHourLogRequestRepository) FindByOrganizationID(ctx context.Context, organizationID int64) (*models.VolunteeringHourLogRequestsResponse, error) {
	response := &models.VolunteeringHourLogRequestsResponse{}

	dbctx := dbctx.Get(ctx)

	db := r.db.
		Model(&models.VolunteeringHourLogRequest{}).
		Limit(dbctx.Limit).
		Where("organization_id = ?", organizationID).
		Count(&response.TotalResults).
		Offset(dbctx.Page * dbctx.Limit)

	if err := db.Find(&response.VolunteeringHourLogRequests).Error; err != nil {
		return response, err
	}

	return response, nil
}

// FindByOpportunityID finds multiple entities by the opportunity ID.
func (r *volunteeringHourLogRequestRepository) FindByOpportunityID(ctx context.Context, opportunityID int64) (*models.VolunteeringHourLogRequestsResponse, error) {
	response := &models.VolunteeringHourLogRequestsResponse{}

	dbctx := dbctx.Get(ctx)

	db := r.db.
		Model(&models.VolunteeringHourLogRequest{}).
		Limit(dbctx.Limit).
		Where("opportunity_id = ?", opportunityID).
		Count(&response.TotalResults).
		Offset(dbctx.Page * dbctx.Limit)

	if err := db.Find(&response.VolunteeringHourLogRequests).Error; err != nil {
		return response, err
	}

	return response, nil
}

// FindByVolunteerID finds multiple entities by the volunteer ID.
func (r *volunteeringHourLogRequestRepository) FindByVolunteerID(ctx context.Context, volunteerID int64) (*models.VolunteeringHourLogRequestsResponse, error) {
	response := &models.VolunteeringHourLogRequestsResponse{}

	dbctx := dbctx.Get(ctx)

	db := r.db.
		Model(&models.VolunteeringHourLogRequest{}).
		Limit(dbctx.Limit).
		Where("volunteer_id = ?", volunteerID).
		Count(&response.TotalResults).
		Offset(dbctx.Page * dbctx.Limit)

	if err := db.Find(&response.VolunteeringHourLogRequests).Error; err != nil {
		return response, err
	}

	return response, nil
}

// Create creates a new User.
func (r *volunteeringHourLogRequestRepository) Create(ctx context.Context, volunteeringHourLogRequest models.VolunteeringHourLogRequest) error {
	return r.db.Create(&volunteeringHourLogRequest).Error
}

// Update updates a User with the ID in the provided User.
func (r *volunteeringHourLogRequestRepository) Update(ctx context.Context, volunteeringHourLogRequest models.VolunteeringHourLogRequest) error {
	return r.db.Model(&models.VolunteeringHourLogRequest{}).Updates(volunteeringHourLogRequest).Error
}

// DeleteByID deletes a User by ID.
func (r *volunteeringHourLogRequestRepository) DeleteByID(ctx context.Context, id int64) error {
	return r.db.Delete(&models.VolunteeringHourLogRequest{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

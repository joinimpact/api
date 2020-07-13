package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// opportunityLimitsRepository stores and controls Opportunities in the database.
type opportunityLimitsRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOpportunityLimitsRepository creates and returns a new OpportunityLimitsRepository.
func NewOpportunityLimitsRepository(db *gorm.DB, logger *zerolog.Logger) models.OpportunityLimitsRepository {
	return &opportunityLimitsRepository{db, logger}
}

// FindByID finds a single entity by ID.
func (r *opportunityLimitsRepository) FindByID(id int64) (*models.OpportunityLimits, error) {
	var opportunityLimits models.OpportunityLimits
	if err := r.db.First(&opportunityLimits, id).Error; err != nil {
		return &opportunityLimits, err
	}
	return &opportunityLimits, nil
}

// FindByOpportunityID finds multiple entities by the organization ID.
func (r *opportunityLimitsRepository) FindByOpportunityID(opportunityID int64) (*models.OpportunityLimits, error) {
	var opportunityLimits models.OpportunityLimits
	if err := r.db.Where("opportunity_id = ?", opportunityID).First(&opportunityLimits).Error; err != nil {
		return &opportunityLimits, err
	}
	return &opportunityLimits, nil
}

// Create creates a new User.
func (r *opportunityLimitsRepository) Create(opportunityLimits models.OpportunityLimits) error {
	return r.db.Create(&opportunityLimits).Error
}

// Update updates a User with the ID in the provided User.
func (r *opportunityLimitsRepository) Update(opportunityLimits models.OpportunityLimits) error {
	return r.db.Model(&models.OpportunityLimits{}).Updates(opportunityLimits).Error
}

// Save saves all fields in the provided entity.
func (r *opportunityLimitsRepository) Save(opportunityLimits models.OpportunityLimits) error {
	return r.db.Save(opportunityLimits).Error
}

// DeleteByID deletes a User by ID.
func (r *opportunityLimitsRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.OpportunityLimits{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

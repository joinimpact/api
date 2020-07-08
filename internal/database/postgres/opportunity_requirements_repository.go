package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// opportunityRequirementsRepository stores and controls OpportunitieRequirements in the database.
type opportunityRequirementsRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOpportunityRequirementsRepository creates and returns a new OpportunityRequirementsRepository.
func NewOpportunityRequirementsRepository(db *gorm.DB, logger *zerolog.Logger) models.OpportunityRequirementsRepository {
	return &opportunityRequirementsRepository{db, logger}
}

// FindByID finds a single entity by ID.
func (r *opportunityRequirementsRepository) FindByID(id int64) (*models.OpportunityRequirements, error) {
	var opportunityRequirements models.OpportunityRequirements
	if err := r.db.First(&opportunityRequirements, id).Error; err != nil {
		return &opportunityRequirements, err
	}
	return &opportunityRequirements, nil
}

// FindByOpportunityID finds multiple entities by the organization ID.
func (r *opportunityRequirementsRepository) FindByOpportunityID(opportunityID int64) (*models.OpportunityRequirements, error) {
	var opportunityRequirements models.OpportunityRequirements
	if err := r.db.Where("opportunity_id = ?", opportunityID).First(&opportunityRequirements).Error; err != nil {
		return &opportunityRequirements, err
	}
	return &opportunityRequirements, nil
}

// Create creates a new User.
func (r *opportunityRequirementsRepository) Create(opportunityRequirements models.OpportunityRequirements) error {
	return r.db.Create(&opportunityRequirements).Error
}

// Update updates a User with the ID in the provided User.
func (r *opportunityRequirementsRepository) Update(opportunityRequirements models.OpportunityRequirements) error {
	return r.db.Model(&models.OpportunityRequirements{}).Updates(opportunityRequirements).Error
}

// DeleteByID deletes a User by ID.
func (r *opportunityRequirementsRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.OpportunityRequirements{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

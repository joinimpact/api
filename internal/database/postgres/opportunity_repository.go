package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// opportunityRepository stores and controls Opportunities in the database.
type opportunityRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOpportunityRepository creates and returns a new OpportunityRepository.
func NewOpportunityRepository(db *gorm.DB, logger *zerolog.Logger) models.OpportunityRepository {
	return &opportunityRepository{db, logger}
}

// FindByID finds a single entity by ID.
func (r *opportunityRepository) FindByID(id int64) (*models.Opportunity, error) {
	var opportunity models.Opportunity
	if err := r.db.First(&opportunity, id).Error; err != nil {
		return &opportunity, err
	}
	return &opportunity, nil
}

// FindByOrganizationID finds multiple entities by the organization ID.
func (r *opportunityRepository) FindByOrganizationID(organizationID int64) ([]models.Opportunity, error) {
	var opportunities []models.Opportunity
	if err := r.db.Where("organization_id = ? AND active = True", organizationID).Find(&opportunities).Error; err != nil {
		return opportunities, err
	}
	return opportunities, nil
}

// FindByCreatorID finds multiple entities by the creator ID.
func (r *opportunityRepository) FindByCreatorID(creatorID int64) ([]models.Opportunity, error) {
	var opportunities []models.Opportunity
	if err := r.db.Where("creator_id = ? AND active = True", creatorID).Find(&opportunities).Error; err != nil {
		return opportunities, err
	}
	return opportunities, nil
}

// Create creates a new User.
func (r *opportunityRepository) Create(opportunity models.Opportunity) error {
	return r.db.Create(&opportunity).Error
}

// Update updates a User with the ID in the provided User.
func (r *opportunityRepository) Update(opportunity models.Opportunity) error {
	return r.db.Model(&models.Opportunity{}).Updates(opportunity).Error
}

// DeleteByID deletes a User by ID.
func (r *opportunityRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.Opportunity{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

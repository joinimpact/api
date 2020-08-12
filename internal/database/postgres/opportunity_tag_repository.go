package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// opportunityTagRepository represents the postgres implementation of the OpportunityTagRepository.
type opportunityTagRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOpportunityTagRepository creates and returns a new OpportunityTagRepository with the provided parameters.
func NewOpportunityTagRepository(db *gorm.DB, logger *zerolog.Logger) models.OpportunityTagRepository {
	return &opportunityTagRepository{
		db,
		logger,
	}
}

// FindByID finds a single entity by ID.
func (r *opportunityTagRepository) FindByID(id int64) (*models.OpportunityTag, error) {
	var opportunityTag models.OpportunityTag
	if err := r.db.Preload("Tag").First(&opportunityTag, id).Error; err != nil {
		return &opportunityTag, err
	}
	return &opportunityTag, nil
}

// FindByOpportunityID finds entities by Opportunity ID.
func (r *opportunityTagRepository) FindByOpportunityID(opportunityID int64) ([]models.OpportunityTag, error) {
	var opportunityTags []models.OpportunityTag
	if err := r.db.Preload("Tag").Where("opportunity_id = ?", opportunityID).Find(&opportunityTags).Error; err != nil {
		return opportunityTags, err
	}
	return opportunityTags, nil
}

// FindOpportunityTagByID finds a single entity by OpportunityID and tag ID.
func (r *opportunityTagRepository) FindOpportunityTagByID(opportunityID int64, tagID int64) (*models.OpportunityTag, error) {
	var opportunityTag models.OpportunityTag
	if err := r.db.Preload("Tag").Where("opportunity_id = ? AND tag_id = ?", opportunityID, tagID).First(&opportunityTag).Error; err != nil {
		return &opportunityTag, err
	}
	return &opportunityTag, nil
}

// Create creates a new entity.
func (r *opportunityTagRepository) Create(opportunityTag models.OpportunityTag) error {
	return r.db.Create(&opportunityTag).Error
}

// Update updates an entity with the ID in the provided entity.
func (r *opportunityTagRepository) Update(opportunityTag models.OpportunityTag) error {
	return r.db.Model(&models.OpportunityTag{}).Updates(opportunityTag).Error
}

// DeleteByID deletes an entity by ID.
func (r *opportunityTagRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.OpportunityTag{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

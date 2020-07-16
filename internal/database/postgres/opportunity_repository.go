package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/pkg/dbctx"
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
func (r *opportunityRepository) FindByID(ctx context.Context, id int64) (*models.Opportunity, error) {
	var opportunity models.Opportunity
	if err := r.db.First(&opportunity, id).Error; err != nil {
		return &opportunity, err
	}
	return &opportunity, nil
}

// FindByOrganizationID finds multiple entities by the organization ID.
func (r *opportunityRepository) FindByOrganizationID(ctx context.Context, organizationID int64) ([]models.Opportunity, error) {
	var opportunities []models.Opportunity
	if err := r.db.
		Limit(dbctx.Get(ctx).Limit).
		Offset(dbctx.Get(ctx).Page*dbctx.Get(ctx).Limit).
		Where("organization_id = ? AND active = True AND LOWER(title) LIKE ?", organizationID, strings.ToLower(fmt.Sprintf("%%%s%%", dbctx.Get(ctx).Query))).
		Find(&opportunities).
		Error; err != nil {
		return opportunities, err
	}
	return opportunities, nil
}

// FindByCreatorID finds multiple entities by the creator ID.
func (r *opportunityRepository) FindByCreatorID(ctx context.Context, creatorID int64) ([]models.Opportunity, error) {
	var opportunities []models.Opportunity
	if err := r.db.Where("creator_id = ? AND active = True", creatorID).Find(&opportunities).Error; err != nil {
		return opportunities, err
	}
	return opportunities, nil
}

// Create creates a new User.
func (r *opportunityRepository) Create(ctx context.Context, opportunity models.Opportunity) error {
	return r.db.Create(&opportunity).Error
}

// Update updates a User with the ID in the provided User.
func (r *opportunityRepository) Update(ctx context.Context, opportunity models.Opportunity) error {
	return r.db.Model(&models.Opportunity{}).Updates(opportunity).Error
}

// Save saves all fields in the provided entity.
func (r *opportunityRepository) Save(ctx context.Context, opportunity models.Opportunity) error {
	return r.db.Save(opportunity).Error
}

// DeleteByID deletes a User by ID.
func (r *opportunityRepository) DeleteByID(ctx context.Context, id int64) error {
	return r.db.Delete(&models.Opportunity{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

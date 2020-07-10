package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// opportunityMembershipRepository stores and controls OpportunityMemberships in the database.
type opportunityMembershipRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOpportunityMembershipRepository creates and returns a new OpportunityMembershipRepository.
func NewOpportunityMembershipRepository(db *gorm.DB, logger *zerolog.Logger) models.OpportunityMembershipRepository {
	return &opportunityMembershipRepository{db, logger}
}

// FindByID finds a single entity by ID.
func (r *opportunityMembershipRepository) FindByID(id int64) (*models.OpportunityMembership, error) {
	var opportunityMembership models.OpportunityMembership
	if err := r.db.First(&opportunityMembership, id).Error; err != nil {
		return &opportunityMembership, err
	}
	return &opportunityMembership, nil
}

// FindByUserID finds multiple entities by the user ID.
func (r *opportunityMembershipRepository) FindByUserID(userID int64) ([]models.OpportunityMembership, error) {
	var opportunityMemberships []models.OpportunityMembership
	if err := r.db.Where("user_id = ? AND active = True", userID).Find(&opportunityMemberships).Error; err != nil {
		return opportunityMemberships, err
	}
	return opportunityMemberships, nil
}

// FindByOpportunityID finds multiple entities by the opportunity ID.
func (r *opportunityMembershipRepository) FindByOpportunityID(opportunityID int64) ([]models.OpportunityMembership, error) {
	var opportunityMemberships []models.OpportunityMembership
	if err := r.db.Where("opportunity_id = ? AND active = True", opportunityID).Find(&opportunityMemberships).Error; err != nil {
		return opportunityMemberships, err
	}
	return opportunityMemberships, nil
}

// FindUserInOpportunity finds a user's membership in a specific opportunity.
func (r *opportunityMembershipRepository) FindUserInOpportunity(opportunityID, userID int64) (*models.OpportunityMembership, error) {
	var opportunityMembership models.OpportunityMembership
	if err := r.db.Where("opportunity_id = ? AND user_id = ? AND active = True", opportunityID, userID).First(&opportunityMembership).Error; err != nil {
		return &opportunityMembership, err
	}
	return &opportunityMembership, nil
}

// Create creates a new User.
func (r *opportunityMembershipRepository) Create(opportunityMembership models.OpportunityMembership) error {
	return r.db.Create(&opportunityMembership).Error
}

// Update updates a User with the ID in the provided User.
func (r *opportunityMembershipRepository) Update(opportunityMembership models.OpportunityMembership) error {
	return r.db.Model(&models.OpportunityMembership{}).Updates(opportunityMembership).Error
}

// DeleteByID deletes a User by ID.
func (r *opportunityMembershipRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.OpportunityMembership{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

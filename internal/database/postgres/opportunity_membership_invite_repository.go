package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// opportunityMembershipInviteRepository represents an implementation of the OpportunityMembershipInviteRepository.
type opportunityMembershipInviteRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOpportunityMembershipInviteRepository creates and returns a new OpportunityMembershipInviteRepository.
func NewOpportunityMembershipInviteRepository(db *gorm.DB, logger *zerolog.Logger) models.OpportunityMembershipInviteRepository {
	return &opportunityMembershipInviteRepository{
		db,
		logger,
	}
}

// FindByID finds a single entity by ID.
func (r *opportunityMembershipInviteRepository) FindByID(id int64) (*models.OpportunityMembershipInvite, error) {
	var opportunityMembershipInvite models.OpportunityMembershipInvite
	if err := r.db.First(&opportunityMembershipInvite, id).Error; err != nil {
		return &opportunityMembershipInvite, err
	}
	return &opportunityMembershipInvite, nil
}

// FindByUserID finds multiple entities by the user ID.
func (r *opportunityMembershipInviteRepository) FindByUserID(userID int64) ([]models.OpportunityMembershipInvite, error) {
	var opportunityMembershipInvites []models.OpportunityMembershipInvite
	if err := r.db.Where("invitee_id = ? AND accepted = False", userID).Find(&opportunityMembershipInvites).Error; err != nil {
		return opportunityMembershipInvites, err
	}
	return opportunityMembershipInvites, nil
}

// FindByUserEmail finds multiple entities by the user Email.
func (r *opportunityMembershipInviteRepository) FindByUserEmail(userEmail string) ([]models.OpportunityMembershipInvite, error) {
	var opportunityMembershipInvites []models.OpportunityMembershipInvite
	if err := r.db.Where("invitee_email = ? AND accepted = False", userEmail).Find(&opportunityMembershipInvites).Error; err != nil {
		return opportunityMembershipInvites, err
	}
	return opportunityMembershipInvites, nil
}

// FindByOpportunityID finds multiple entities by the opportunity ID.
func (r *opportunityMembershipInviteRepository) FindByOpportunityID(opportunityID int64) ([]models.OpportunityMembershipInvite, error) {
	var opportunityMembershipInvites []models.OpportunityMembershipInvite
	if err := r.db.Where("opportunity_id = ? AND accepted = False", opportunityID).Find(&opportunityMembershipInvites).Error; err != nil {
		return opportunityMembershipInvites, err
	}
	return opportunityMembershipInvites, nil
}

// FindByOpportunityIDs finds multiple entities by multiple opportunity IDs.
func (r *opportunityMembershipInviteRepository) FindByOpportunityIDs(ids []int64) ([]models.OpportunityMembershipInvite, error) {
	var opportunityMembershipInvites []models.OpportunityMembershipInvite
	if err := r.db.Preload("Opportunity").Where("opportunity_id in (?) AND accepted = False", ids).Find(&opportunityMembershipInvites).Error; err != nil {
		return opportunityMembershipInvites, err
	}
	return opportunityMembershipInvites, nil
}

// FindInOpportunityByUserID finds a membership invite in an opportunity by user ID.
func (r *opportunityMembershipInviteRepository) FindInOpportunityByUserID(opportunityID, userID int64) (*models.OpportunityMembershipInvite, error) {
	var opportunityMembershipInvite models.OpportunityMembershipInvite
	if err := r.db.Where("opportunity_id = ? AND invitee_id = ? AND accepted = False", opportunityID, userID).First(&opportunityMembershipInvite).Error; err != nil {
		return &opportunityMembershipInvite, err
	}
	return &opportunityMembershipInvite, nil
}

// FindInOpportunityByEmail finds a membership invite in an opportunity by user email.
func (r *opportunityMembershipInviteRepository) FindInOpportunityByEmail(opportunityID int64, email string) (*models.OpportunityMembershipInvite, error) {
	var opportunityMembershipInvite models.OpportunityMembershipInvite
	if err := r.db.Where("opportunity_id = ? AND invitee_email = ? AND accepted = False", opportunityID, email).First(&opportunityMembershipInvite).Error; err != nil {
		return &opportunityMembershipInvite, err
	}
	return &opportunityMembershipInvite, nil
}

// Create creates a new entity.
func (r *opportunityMembershipInviteRepository) Create(opportunityMembershipInvite models.OpportunityMembershipInvite) error {
	return r.db.Create(&opportunityMembershipInvite).Error
}

// Update updates an entity with the ID in the provided entity.
func (r *opportunityMembershipInviteRepository) Update(opportunityMembershipInvite models.OpportunityMembershipInvite) error {
	return r.db.Model(&models.OpportunityMembershipInvite{}).Updates(opportunityMembershipInvite).Error
}

// DeleteByID deletes an entity by ID.
func (r *opportunityMembershipInviteRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.OpportunityMembershipInvite{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

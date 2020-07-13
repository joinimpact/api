package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// opportunityMembershipRequestRepository represents an implementation of the OpportunityMembershipRequestRepository.
type opportunityMembershipRequestRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOpportunityMembershipRequestRepository creates and returns a new OpportunityMembershipRequestRepository.
func NewOpportunityMembershipRequestRepository(db *gorm.DB, logger *zerolog.Logger) models.OpportunityMembershipRequestRepository {
	return &opportunityMembershipRequestRepository{
		db,
		logger,
	}
}

// FindByID finds a single entity by ID.
func (r *opportunityMembershipRequestRepository) FindByID(id int64) (*models.OpportunityMembershipRequest, error) {
	var opportunityMembershipRequest models.OpportunityMembershipRequest
	if err := r.db.First(&opportunityMembershipRequest, id).Error; err != nil {
		return &opportunityMembershipRequest, err
	}
	return &opportunityMembershipRequest, nil
}

// FindByVolunteerID finds multiple entities by the volunteer ID.
func (r *opportunityMembershipRequestRepository) FindByVolunteerID(volunteerID int64) ([]models.OpportunityMembershipRequest, error) {
	var opportunityMembershipRequests []models.OpportunityMembershipRequest
	if err := r.db.Where("volunteer_id = ? AND accepted = False", volunteerID).Find(&opportunityMembershipRequests).Error; err != nil {
		return opportunityMembershipRequests, err
	}
	return opportunityMembershipRequests, nil
}

// FindByOpportunityID finds multiple entities by the opportunity ID.
func (r *opportunityMembershipRequestRepository) FindByOpportunityID(opportunityID int64) ([]models.OpportunityMembershipRequest, error) {
	var opportunityMembershipRequests []models.OpportunityMembershipRequest
	if err := r.db.Where("opportunity_id = ? AND accepted = False", opportunityID).Find(&opportunityMembershipRequests).Error; err != nil {
		return opportunityMembershipRequests, err
	}
	return opportunityMembershipRequests, nil
}

// FindInOpportunityByVolunteerID finds a single entity by opportunity and volunteer ID.
func (r *opportunityMembershipRequestRepository) FindInOpportunityByVolunteerID(opportunityID, volunteerID int64) (*models.OpportunityMembershipRequest, error) {
	var opportunityMembership models.OpportunityMembershipRequest
	if err := r.db.Where("opportunity_id = ? AND volunteer_id = ?", opportunityID, volunteerID).First(&opportunityMembership).Error; err != nil {
		return &opportunityMembership, err
	}
	return &opportunityMembership, nil
}

// Create creates a new entity.
func (r *opportunityMembershipRequestRepository) Create(opportunityMembershipRequest models.OpportunityMembershipRequest) error {
	return r.db.Create(&opportunityMembershipRequest).Error
}

// Update updates an entity with the ID in the provided entity.
func (r *opportunityMembershipRequestRepository) Update(opportunityMembershipRequest models.OpportunityMembershipRequest) error {
	return r.db.Model(&models.OpportunityMembershipRequest{}).Updates(opportunityMembershipRequest).Error
}

// DeleteByID deletes an entity by ID.
func (r *opportunityMembershipRequestRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.OpportunityMembershipRequest{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

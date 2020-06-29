package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// organizationMembershipInviteRepository represents an implementation of the OrganizationMembershipInviteRepository.
type organizationMembershipInviteRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOrganizationMembershipInviteRepository creates and returns a new OrganizationMembershipInviteRepository.
func NewOrganizationMembershipInviteRepository(db *gorm.DB, logger *zerolog.Logger) models.OrganizationMembershipInviteRepository {
	return &organizationMembershipInviteRepository{
		db,
		logger,
	}
}

// FindByID finds a single entity by ID.
func (r *organizationMembershipInviteRepository) FindByID(id int64) (*models.OrganizationMembershipInvite, error) {
	var organizationMembershipInvite models.OrganizationMembershipInvite
	if err := r.db.First(&organizationMembershipInvite, id).Error; err != nil {
		return &organizationMembershipInvite, err
	}
	return &organizationMembershipInvite, nil
}

// FindByUserID finds multiple entities by the user ID.
func (r *organizationMembershipInviteRepository) FindByUserID(userID int64) ([]models.OrganizationMembershipInvite, error) {
	var organizationMembershipInvites []models.OrganizationMembershipInvite
	if err := r.db.Where("invitee_id = ? AND accepted = False", userID).Find(&organizationMembershipInvites).Error; err != nil {
		return organizationMembershipInvites, err
	}
	return organizationMembershipInvites, nil
}

// FindByUserEmail finds multiple entities by the user Email.
func (r *organizationMembershipInviteRepository) FindByUserEmail(userEmail string) ([]models.OrganizationMembershipInvite, error) {
	var organizationMembershipInvites []models.OrganizationMembershipInvite
	if err := r.db.Where("invitee_email = ? AND accepted = False", userEmail).Find(&organizationMembershipInvites).Error; err != nil {
		return organizationMembershipInvites, err
	}
	return organizationMembershipInvites, nil
}

// FindByOrganizationID finds multiple entities by the organization ID.
func (r *organizationMembershipInviteRepository) FindByOrganizationID(organizationID int64) ([]models.OrganizationMembershipInvite, error) {
	var organizationMembershipInvites []models.OrganizationMembershipInvite
	if err := r.db.Where("organization_id = ? AND accepted = False", organizationID).Find(&organizationMembershipInvites).Error; err != nil {
		return organizationMembershipInvites, err
	}
	return organizationMembershipInvites, nil
}

// Create creates a new entity.
func (r *organizationMembershipInviteRepository) Create(organizationMembershipInvite models.OrganizationMembershipInvite) error {
	return r.db.Create(&organizationMembershipInvite).Error
}

// Update updates an entity with the ID in the provided entity.
func (r *organizationMembershipInviteRepository) Update(organizationMembershipInvite models.OrganizationMembershipInvite) error {
	return r.db.Model(&models.OrganizationMembershipInvite{}).Updates(organizationMembershipInvite).Error
}

// DeleteByID deletes an entity by ID.
func (r *organizationMembershipInviteRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.OrganizationMembershipInvite{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// organizationMembershipRepository stores and controls OrganizationMemberships in the database.
type organizationMembershipRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOrganizationMembershipRepository creates and returns a new OrganizationMembershipRepository.
func NewOrganizationMembershipRepository(db *gorm.DB, logger *zerolog.Logger) models.OrganizationMembershipRepository {
	return &organizationMembershipRepository{db, logger}
}

// FindByID finds a single entity by ID.
func (r *organizationMembershipRepository) FindByID(id int64) (*models.OrganizationMembership, error) {
	var organizationMembership models.OrganizationMembership
	if err := r.db.First(&organizationMembership, id).Error; err != nil {
		return &organizationMembership, err
	}
	return &organizationMembership, nil
}

// FindByUserID finds multiple entities by the user ID.
func (r *organizationMembershipRepository) FindByUserID(userID int64) ([]models.OrganizationMembership, error) {
	var organizationMemberships []models.OrganizationMembership
	if err := r.db.Where("user_id = ? AND active = True", userID).Find(&organizationMemberships).Error; err != nil {
		return organizationMemberships, err
	}
	return organizationMemberships, nil
}

// FindByOrganizationID finds multiple entities by the organization ID.
func (r *organizationMembershipRepository) FindByOrganizationID(organizationID int64) ([]models.OrganizationMembership, error) {
	var organizationMemberships []models.OrganizationMembership
	if err := r.db.Where("organization_id = ? AND active = True", organizationID).Find(&organizationMemberships).Error; err != nil {
		return organizationMemberships, err
	}
	return organizationMemberships, nil
}

// FindUserInOrganization finds a user's membership in a specific organization.
func (r *organizationMembershipRepository) FindUserInOrganization(organizationID, userID int64) (*models.OrganizationMembership, error) {
	var organizationMembership models.OrganizationMembership
	if err := r.db.Where("organization_id = ? AND user_id = ? AND active = True", organizationID, userID).First(&organizationMembership).Error; err != nil {
		return &organizationMembership, err
	}
	return &organizationMembership, nil
}

// Create creates a new User.
func (r *organizationMembershipRepository) Create(organizationMembership models.OrganizationMembership) error {
	return r.db.Create(&organizationMembership).Error
}

// Update updates a User with the ID in the provided User.
func (r *organizationMembershipRepository) Update(organizationMembership models.OrganizationMembership) error {
	return r.db.Model(&models.OrganizationMembership{}).Updates(organizationMembership).Error
}

// DeleteByID deletes a User by ID.
func (r *organizationMembershipRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.OrganizationMembership{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// organizationRepository stores and controls Organizations in the database.
type organizationRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOrganizationRepository creates and returns a new OrganizationRepository.
func NewOrganizationRepository(db *gorm.DB, logger *zerolog.Logger) models.OrganizationRepository {
	return &organizationRepository{db, logger}
}

// FindByID finds a single User by ID.
func (r *organizationRepository) FindByID(id int64) (*models.Organization, error) {
	var organization models.Organization
	if err := r.db.First(&organization, id).Error; err != nil || !organization.Active {
		return &organization, err
	}
	return &organization, nil
}

// FindByCreatorID finds multiple entities by the creator's ID.
func (r *organizationRepository) FindByCreatorID(creatorID int64) ([]models.Organization, error) {
	var organizations []models.Organization
	if err := r.db.Where("creator_id = ? AND active = True", creatorID).Find(&organizations).Error; err != nil {
		return organizations, err
	}
	return organizations, nil
}

// Create creates a new User.
func (r *organizationRepository) Create(organization models.Organization) error {
	return r.db.Create(&organization).Error
}

// Update updates a User with the ID in the provided User.
func (r *organizationRepository) Update(organization models.Organization) error {
	return r.db.Model(&models.Organization{}).Updates(organization).Error
}

// DeleteByID deletes a User by ID.
func (r *organizationRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.Organization{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

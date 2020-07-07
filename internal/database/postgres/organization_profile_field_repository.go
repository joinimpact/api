package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// organizationProfileFieldRepository is the postgres implementation of OrganizationProfileFieldRepository.
type organizationProfileFieldRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOrganizationProfileFieldRepository creates and returns a new OrganizationProfileFieldRepository.
func NewOrganizationProfileFieldRepository(db *gorm.DB, logger *zerolog.Logger) models.OrganizationProfileFieldRepository {
	return &organizationProfileFieldRepository{
		db,
		logger,
	}
}

// FindByID finds a single entity by ID.
func (r *organizationProfileFieldRepository) FindByID(id int64) (*models.OrganizationProfileField, error) {
	var profileField models.OrganizationProfileField
	if err := r.db.First(&profileField, id).Error; err != nil {
		return &profileField, err
	}
	return &profileField, nil
}

// FindByOrganizationID finds entities by OrganizationID.
func (r *organizationProfileFieldRepository) FindByOrganizationID(id int64) ([]models.OrganizationProfileField, error) {
	var profileFields []models.OrganizationProfileField
	if err := r.db.Where("organization_id = ?", id).Find(&profileFields).Error; err != nil {
		return profileFields, err
	}
	return profileFields, nil
}

// FindOrganizationFieldByName finds a single entity by UserID and field name.
func (r *organizationProfileFieldRepository) FindOrganizationFieldByName(id int64, name string) (*models.OrganizationProfileField, error) {
	var profileField models.OrganizationProfileField
	if err := r.db.Where("organization_id = ? AND name = ?", id, name).Find(&profileField).Error; err != nil {
		return &profileField, err
	}
	return &profileField, nil
}

// Create creates a new entity.
func (r *organizationProfileFieldRepository) Create(profileField models.OrganizationProfileField) error {
	return r.db.Create(&profileField).Error
}

// Update updates an entity with the ID in the provided entity.
func (r *organizationProfileFieldRepository) Update(profileField models.OrganizationProfileField) error {
	return r.db.Model(&models.OrganizationProfileField{}).Updates(&profileField).Error
}

// DeleteByID deletes an entity by ID.
func (r *organizationProfileFieldRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.OrganizationProfileField{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

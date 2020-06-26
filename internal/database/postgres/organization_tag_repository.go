package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// organizationTagRepository represents the postgres implementation of the OrganizationTagRepository.
type organizationTagRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOrganizationTagRepository creates and returns a new OrganizationTagRepository with the provided parameters.
func NewOrganizationTagRepository(db *gorm.DB, logger *zerolog.Logger) models.OrganizationTagRepository {
	return &organizationTagRepository{
		db,
		logger,
	}
}

// FindByID finds a single entity by ID.
func (r *organizationTagRepository) FindByID(id int64) (*models.OrganizationTag, error) {
	var organizationTag models.OrganizationTag
	if err := r.db.First(&organizationTag, id).Error; err != nil {
		return &organizationTag, err
	}
	return &organizationTag, nil
}

// FindByOrganizationID finds entities by Organization ID.
func (r *organizationTagRepository) FindByOrganizationID(organizationID int64) ([]models.OrganizationTag, error) {
	var organizationTags []models.OrganizationTag
	if err := r.db.Where("organization_id = ?", organizationID).Find(&organizationTags).Error; err != nil {
		return organizationTags, err
	}
	return organizationTags, nil
}

// FindOrganizationTagByID finds a single entity by OrganizationID and tag ID.
func (r *organizationTagRepository) FindOrganizationTagByID(organizationID int64, tagID int64) (*models.OrganizationTag, error) {
	var organizationTag models.OrganizationTag
	if err := r.db.Where("organization_id = ? AND tag_id = ?", organizationID, tagID).First(&organizationTag).Error; err != nil {
		return &organizationTag, err
	}
	return &organizationTag, nil
}

// Create creates a new entity.
func (r *organizationTagRepository) Create(organizationTag models.OrganizationTag) error {
	return r.db.Create(&organizationTag).Error
}

// Update updates an entity with the ID in the provided entity.
func (r *organizationTagRepository) Update(organizationTag models.OrganizationTag) error {
	return r.db.Model(&models.OrganizationTag{}).Updates(organizationTag).Error
}

// DeleteByID deletes an entity by ID.
func (r *organizationTagRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.OrganizationTag{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

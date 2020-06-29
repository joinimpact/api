package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// tagRepository represents the postgres implementation of the TagRepository.
type tagRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewTagRepository creates and returns a new TagRepository with the provided parameters.
func NewTagRepository(db *gorm.DB, logger *zerolog.Logger) models.TagRepository {
	return &tagRepository{
		db,
		logger,
	}
}

// FindByID finds a single entity by ID.
func (r *tagRepository) FindByID(id int64) (*models.Tag, error) {
	var tag models.Tag
	if err := r.db.First(&tag, id).Error; err != nil {
		return &tag, err
	}
	return &tag, nil
}

// FindByCategory finds entities by Category.
func (r *tagRepository) FindByCategory(category int) ([]models.Tag, error) {
	var tags []models.Tag
	if err := r.db.Where("category = ?", category).Find(&tags).Error; err != nil {
		return tags, err
	}
	return tags, nil
}

// FindByName finds a single entity by name.
func (r *tagRepository) FindByName(name string) (*models.Tag, error) {
	var tag models.Tag
	if err := r.db.Where("name = ?", name).Find(&tag).Error; err != nil {
		return &tag, err
	}
	return &tag, nil
}

// SearchTags searches for tags with a query string.
func (r *tagRepository) SearchTags(query string, limit int) ([]models.Tag, error) {
	var tags []models.Tag
	if err := r.db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query)).Limit(limit).Find(&tags).Error; err != nil {
		return tags, err
	}
	return tags, nil
}

// Create creates a new entity.
func (r *tagRepository) Create(tag models.Tag) error {
	return r.db.Create(&tag).Error
}

// Update updates an entity with the ID in the provided entity.
func (r *tagRepository) Update(tag models.Tag) error {
	return r.db.Model(&models.Tag{}).Updates(tag).Error
}

// DeleteByID deletes an entity by ID.
func (r *tagRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.Tag{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

package migrations

import "github.com/jinzhu/gorm"

// MigrationService provides methods for migrating a Postgres database with the models.
type MigrationService struct {
	db *gorm.DB
}

// NewMigrationService creates and returns a new MigrationService.
func NewMigrationService(db *gorm.DB) *MigrationService {
	return &MigrationService{db}
}

// Migrate takes all database models passed in and runs gorm's auto migrate
// function.
func (s *MigrationService) Migrate(models ...interface{}) error {
	s.db.AutoMigrate(models...)
	return nil
}

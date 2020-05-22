package migrations

import "github.com/jinzhu/gorm"

// MigrationService provides methods for migrating a Postgres database with the models.
type MigrationService struct {
	db *gorm.DB
}

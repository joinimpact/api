package core

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/config"
	"github.com/rs/zerolog"
)

// DatabaseService provides functions for connecting to a database.
type DatabaseService struct {
	config *config.Config
	logger *zerolog.Logger
}

// NewDatabaseService creates and returns a new DatabaseService.
func NewDatabaseService(config *config.Config, logger *zerolog.Logger) *DatabaseService {
	return &DatabaseService{config, logger}
}

// DatabaseConnect attempts to establish a connection with the database and
// returns a gorm.DB pointer.
func (s *DatabaseService) DatabaseConnect() (*gorm.DB, error) {
	return gorm.Open("postgres", fmt.Sprintf(
		"sslmode=disable host=%s port=%d user=%s dbname=%s password=%s",
		s.config.DatabaseHost,
		s.config.DatabasePort,
		s.config.DatabaseUser,
		s.config.DatabaseName,
		s.config.DatabasePassword),
	)
}

package postgres

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/pkg/dbctx"
	"github.com/rs/zerolog"
)

// opportunityApplicationQuestionRepositoryg stores and controls messages in the database.
type opportunityApplicationQuestionRepositoryg struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOpportunityApplicationQuestionRepository creates and returns a new OpportunityApplicationQuestionRepository.
func NewOpportunityApplicationQuestionRepository(db *gorm.DB, logger *zerolog.Logger) models.OpportunityApplicationQuestionRepository {
	return &opportunityApplicationQuestionRepositoryg{db, logger}
}

// FindByID finds a single entity by ID.
func (r *opportunityApplicationQuestionRepositoryg) FindByID(ctx context.Context, id int64) (*models.OpportunityApplicationQuestion, error) {
	var opportunityApplicationQuestion models.OpportunityApplicationQuestion
	if err := r.db.First(&opportunityApplicationQuestion, id).Error; err != nil {
		return &opportunityApplicationQuestion, err
	}
	return &opportunityApplicationQuestion, nil
}

// FindByOpportunityID finds multiple entities by the opportunity ID.
func (r *opportunityApplicationQuestionRepositoryg) FindByOpportunityID(ctx context.Context, opportunityID int64) (*models.OpportunityApplicationQuestionsResponse, error) {
	opportunityApplicationQuestionsResponse := &models.OpportunityApplicationQuestionsResponse{}

	dbctx := dbctx.Get(ctx)

	db := r.db.
		Model(&models.OpportunityApplicationQuestion{}).
		Limit(dbctx.Limit).
		Where("opportunity_id = ?", opportunityID).
		Count(&opportunityApplicationQuestionsResponse.TotalResults).
		Offset(dbctx.Page * dbctx.Limit)

	if err := db.Find(&opportunityApplicationQuestionsResponse.OpportunityApplicationQuestions).Error; err != nil {
		return opportunityApplicationQuestionsResponse, err
	}
	return opportunityApplicationQuestionsResponse, nil
}

// Create creates a new User.
func (r *opportunityApplicationQuestionRepositoryg) Create(ctx context.Context, opportunityApplicationQuestion models.OpportunityApplicationQuestion) error {
	return r.db.Create(&opportunityApplicationQuestion).Error
}

// Update updates a User with the ID in the provided User.
func (r *opportunityApplicationQuestionRepositoryg) Update(ctx context.Context, opportunityApplicationQuestion models.OpportunityApplicationQuestion) error {
	return r.db.Model(&models.OpportunityApplicationQuestion{}).Updates(opportunityApplicationQuestion).Error
}

// DeleteByID deletes a User by ID.
func (r *opportunityApplicationQuestionRepositoryg) DeleteByID(ctx context.Context, id int64) error {
	return r.db.Delete(&models.Message{
		Model: models.Model{
			ID: id,
		},
	}).Error
}

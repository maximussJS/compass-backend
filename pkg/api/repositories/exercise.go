package repositories

import (
	"compass-backend/pkg/api/models"
	fx_utils "compass-backend/pkg/common/fx"
	common_infrastructure "compass-backend/pkg/common/infrastructure"
	"context"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type IExerciseRepository interface {
	Create(ctx context.Context, exercise models.Exercise) (string, error)
	GetById(ctx context.Context, id string) (*models.Exercise, error)
	GetByNameAndCreatorId(ctx context.Context, name, creatorId string) (*models.Exercise, error)
	GetAllByCreatorId(ctx context.Context, creatorId string) ([]models.Exercise, error)
	Delete(ctx context.Context, id string) error
}

type exerciseRepositoryParams struct {
	fx.In

	Database common_infrastructure.IDatabase
}

type exerciseRepository struct {
	db *gorm.DB
}

func FxExerciseRepository() fx.Option {
	return fx_utils.AsProvider(newExerciseRepository, new(IExerciseRepository))
}

func newExerciseRepository(params exerciseRepositoryParams) IExerciseRepository {
	return &exerciseRepository{
		db: params.Database.GetInstance(),
	}
}

func (r *exerciseRepository) Create(ctx context.Context, exercise models.Exercise) (string, error) {
	err := r.db.WithContext(ctx).Create(&exercise).Error
	if err != nil {
		return "", err
	}

	return exercise.Id, nil
}

func (r *exerciseRepository) GetById(ctx context.Context, id string) (*models.Exercise, error) {
	var exercise *models.Exercise
	err := r.db.WithContext(ctx).Where("id = ?", id).First(exercise).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, nil
		}

		return nil, err
	}

	return exercise, nil
}

func (r *exerciseRepository) GetByNameAndCreatorId(ctx context.Context, name, creatorId string) (*models.Exercise, error) {
	var exercise *models.Exercise
	err := r.db.WithContext(ctx).Where("name = ? AND creator_id = ?", name, creatorId).First(exercise).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, nil
		}

		return nil, err
	}

	return exercise, nil
}

func (r *exerciseRepository) GetAllByCreatorId(ctx context.Context, creatorId string) ([]models.Exercise, error) {
	var exercises []models.Exercise
	err := r.db.WithContext(ctx).Where("creator_id = ?", creatorId).Find(&exercises).Error
	if err != nil {
		return nil, err
	}

	return exercises, nil
}

func (r *exerciseRepository) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Exercise{}).Error
	if err != nil {
		return err
	}

	return nil
}

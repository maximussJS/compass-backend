package repositories

import (
	"compass-backend/pkg/api/models"
	fx_utils "compass-backend/pkg/common/fx"
	common_infrastructure "compass-backend/pkg/common/infrastructure"
	"context"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type IExerciseMediaRepository interface {
	Create(ctx context.Context, exercise models.ExerciseMedia) (string, error)
	GetById(ctx context.Context, id string) (*models.ExerciseMedia, error)
	Delete(ctx context.Context, id string) error
}

type exerciseMediaRepositoryParams struct {
	fx.In

	Database common_infrastructure.IDatabase
}

type exerciseMediaRepository struct {
	db *gorm.DB
}

func FxExerciseMediaRepository() fx.Option {
	return fx_utils.AsProvider(newExerciseMediaRepository, new(IExerciseMediaRepository))
}

func newExerciseMediaRepository(params exerciseMediaRepositoryParams) IExerciseMediaRepository {
	return &exerciseMediaRepository{
		db: params.Database.GetInstance(),
	}
}

func (r *exerciseMediaRepository) Create(ctx context.Context, exercise models.ExerciseMedia) (string, error) {
	err := r.db.WithContext(ctx).Create(&exercise).Error
	if err != nil {
		return "", err
	}

	return exercise.Id, nil
}

func (r *exerciseMediaRepository) GetById(ctx context.Context, id string) (*models.ExerciseMedia, error) {
	var exercise *models.ExerciseMedia
	err := r.db.WithContext(ctx).Where("id = ?", id).First(exercise).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, nil
		}

		return nil, err
	}

	return exercise, nil
}

func (r *exerciseMediaRepository) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.ExerciseMedia{}).Error
	if err != nil {
		return err
	}

	return nil
}

package repositories

import (
	"compass-backend/pkg/api/models"
	fx_utils "compass-backend/pkg/common/fx"
	gorm_utils "compass-backend/pkg/common/gorm"
	common_infrastructure "compass-backend/pkg/common/infrastructure"
	"context"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type ICategoryRepository interface {
	Create(ctx context.Context, category models.Category) (uint, error)
	GetById(ctx context.Context, id uint) (*models.Category, error)
	GetByName(ctx context.Context, name string) (*models.Category, error)
	List(ctx context.Context, limit, offset int) ([]models.Category, error)
	UpdateById(ctx context.Context, id uint, category models.Category) error
	DeleteById(ctx context.Context, id uint) error
}

type categoryRepositoryParams struct {
	fx.In

	Database common_infrastructure.IDatabase
}

type categoryRepository struct {
	db *gorm.DB
}

func FxCategoryRepository() fx.Option {
	return fx_utils.AsProvider(newCategoryRepository, new(ICategoryRepository))
}

func newCategoryRepository(params categoryRepositoryParams) ICategoryRepository {
	return &categoryRepository{
		db: params.Database.GetInstance(),
	}
}

func (r *categoryRepository) Create(ctx context.Context, category models.Category) (uint, error) {
	err := r.db.WithContext(ctx).Create(&category).Error
	if err != nil {
		return 0, err
	}

	return category.Id, nil
}

func (r *categoryRepository) GetById(ctx context.Context, id uint) (*models.Category, error) {
	category := &models.Category{}
	err := r.db.WithContext(ctx).First(category, id).Error
	if err != nil {
		if gorm_utils.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return category, nil
}

func (r *categoryRepository) GetByName(ctx context.Context, name string) (*models.Category, error) {
	category := &models.Category{}
	err := r.db.WithContext(ctx).Where("name = ?", name).First(category).Error
	if err != nil {
		if gorm_utils.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return category, nil
}

func (r *categoryRepository) List(ctx context.Context, limit, offset int) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("id ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) UpdateById(ctx context.Context, id uint, category models.Category) error {
	err := r.db.WithContext(ctx).Model(&models.Category{}).Where("id = ?", id).Updates(category).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepository) DeleteById(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&models.Category{}, id).Error
	if err != nil {
		return err
	}

	return nil
}

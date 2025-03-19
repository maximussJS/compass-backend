package repositories

import (
	fx_utils "compass-backend/pkg/common/fx"
	gorm_utils "compass-backend/pkg/common/gorm"
	"compass-backend/pkg/common/infrastructure"
	"compass-backend/pkg/common/models"
	"context"
	"fmt"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(ctx context.Context, user models.User) (string, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetById(ctx context.Context, id string) (*models.User, error)
	UpdateById(ctx context.Context, id string, user models.User) error
}

type userRepositoryParams struct {
	fx.In

	Database infrastructure.IDatabase
}

type userRepository struct {
	db *gorm.DB
}

func FxUserRepository() fx.Option {
	return fx_utils.AsProvider(newUserRepository, new(IUserRepository))
}

func newUserRepository(params userRepositoryParams) IUserRepository {
	return &userRepository{
		db: params.Database.GetInstance(),
	}
}

func (r *userRepository) Create(ctx context.Context, user models.User) (string, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return user.Id, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}

	err := r.db.WithContext(ctx).Preload("Teams").Where("email = ?", email).First(user).Error

	if err != nil {
		if gorm_utils.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetById(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}

	err := r.db.WithContext(ctx).Preload("Teams").Where("id = ?", id).First(user).Error

	if err != nil {
		if gorm_utils.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

func (r *userRepository) UpdateById(ctx context.Context, id string, user models.User) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Updates(&user).Error
	if err != nil {
		return fmt.Errorf("failed to update user by id: %w", err)
	}

	return nil
}

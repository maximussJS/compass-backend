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

type IAccessKeyRepository interface {
	Create(ctx context.Context, key models.AccessKey) (string, error)
	GetById(ctx context.Context, id string) (*models.AccessKey, error)
	GetByIdAndUserId(ctx context.Context, id, userId string) (*models.AccessKey, error)
	GetByNameAndTeamId(ctx context.Context, name, keyId string) (*models.AccessKey, error)
	GetByTeamId(ctx context.Context, keyId string, limit, offset int) ([]models.AccessKey, error)
	GetByIdAndTeamId(ctx context.Context, id, keyId string) (*models.AccessKey, error)
	UpdateById(ctx context.Context, id string, key models.AccessKey) error
	DeleteById(ctx context.Context, id string) error
}

type accessKeyRepositoryParams struct {
	fx.In

	Database infrastructure.IDatabase
}

type accessKeyRepository struct {
	db *gorm.DB
}

func FxAccessKeyRepository() fx.Option {
	return fx_utils.AsProvider(newAccessKeyRepository, new(IAccessKeyRepository))
}

func newAccessKeyRepository(params accessKeyRepositoryParams) IAccessKeyRepository {
	return &accessKeyRepository{
		db: params.Database.GetInstance(),
	}
}

func (r *accessKeyRepository) Create(ctx context.Context, key models.AccessKey) (string, error) {
	err := r.db.WithContext(ctx).Create(&key).Error
	if err != nil {
		return "", fmt.Errorf("failed to create access key: %w", err)
	}

	return key.Id, nil
}

func (r *accessKeyRepository) GetById(ctx context.Context, id string) (*models.AccessKey, error) {
	accessKey := &models.AccessKey{}
	err := r.db.WithContext(ctx).Where("id = ?", id).First(accessKey).Error
	if err != nil {
		if gorm_utils.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get access key by id: %w", err)
	}

	return accessKey, nil
}

func (r *accessKeyRepository) GetByIdAndUserId(ctx context.Context, id, userId string) (*models.AccessKey, error) {
	accessKey := &models.AccessKey{}
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userId).First(accessKey).Error
	if err != nil {
		if gorm_utils.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get access key by id and user id: %w", err)
	}

	return accessKey, nil
}

func (r *accessKeyRepository) GetByNameAndTeamId(ctx context.Context, name, keyId string) (*models.AccessKey, error) {
	accessKey := &models.AccessKey{}
	err := r.db.WithContext(ctx).Where("name = ? AND team_id = ?", name, keyId).First(accessKey).Error
	if err != nil {
		if gorm_utils.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get access key by name and team id: %w", err)
	}

	return accessKey, nil
}

func (r *accessKeyRepository) GetByTeamId(ctx context.Context, keyId string, limit, offset int) ([]models.AccessKey, error) {
	var accessKeys []models.AccessKey
	err := r.db.WithContext(ctx).Where("team_id = ?", keyId).Limit(limit).Offset(offset).Find(&accessKeys).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get access keys by team id: %w", err)
	}

	return accessKeys, nil
}

func (r *accessKeyRepository) GetByIdAndTeamId(ctx context.Context, id, keyId string) (*models.AccessKey, error) {
	accessKey := &models.AccessKey{}
	err := r.db.WithContext(ctx).Where("id = ? AND team_id = ?", id, keyId).First(accessKey).Error
	if err != nil {
		if gorm_utils.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get access key by id and team id: %w", err)
	}

	return accessKey, nil
}

func (r *accessKeyRepository) UpdateById(ctx context.Context, id string, key models.AccessKey) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Updates(&key).Error
	if err != nil {
		return fmt.Errorf("failed to update access key: %w", err)
	}

	return nil
}

func (r *accessKeyRepository) DeleteById(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.AccessKey{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete access key: %w", err)
	}

	return nil
}

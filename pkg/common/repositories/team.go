package repositories

import (
	fx_utils "compass-backend/pkg/common/fx"
	gorm_utils "compass-backend/pkg/common/gorm"
	"compass-backend/pkg/common/infrastructure"
	"compass-backend/pkg/common/models"
	"context"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ITeamRepository interface {
	Create(ctx context.Context, team models.Team) (string, error)
	GetById(ctx context.Context, id string) (*models.Team, error)
	GetByOwnerId(ctx context.Context, ownerId string) (*models.Team, error)
	UpdateById(ctx context.Context, id string, team models.Team) error
	DeleteById(ctx context.Context, id string) error
}

type teamRepositoryParams struct {
	fx.In

	Database infrastructure.IDatabase
}

type teamRepository struct {
	db *gorm.DB
}

func FxTeamRepository() fx.Option {
	return fx_utils.AsProvider(newTeamRepository, new(ITeamRepository))
}

func newTeamRepository(params teamRepositoryParams) ITeamRepository {
	return &teamRepository{
		db: params.Database.GetInstance(),
	}
}

func (r *teamRepository) Create(ctx context.Context, team models.Team) (string, error) {
	err := r.db.WithContext(ctx).Create(&team).Error
	if err != nil {
		return "", err
	}

	return team.Id, nil
}

func (r *teamRepository) GetById(ctx context.Context, id string) (*models.Team, error) {
	team := &models.Team{}
	err := r.db.WithContext(ctx).Preload(clause.Associations).Where("id = ?", id).First(team).Error
	if err != nil {
		if gorm_utils.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return team, nil
}

func (r *teamRepository) GetByOwnerId(ctx context.Context, ownerId string) (*models.Team, error) {
	team := &models.Team{}
	err := r.db.WithContext(ctx).Where("owner_id = ?", ownerId).First(team).Error
	if err != nil {
		if gorm_utils.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return team, nil
}

func (r *teamRepository) UpdateById(ctx context.Context, id string, team models.Team) error {
	err := r.db.WithContext(ctx).Model(&models.Team{}).Where("id = ?", id).Updates(team).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *teamRepository) DeleteById(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Delete(&models.Team{}, id).Error
	if err != nil {
		return err
	}

	return nil
}

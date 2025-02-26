package repositories

import (
	"compass-backend/pkg/common/constants"
	fx_utils "compass-backend/pkg/common/fx"
	gorm_utils "compass-backend/pkg/common/gorm"
	"compass-backend/pkg/common/infrastructure"
	"compass-backend/pkg/common/models"
	"context"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type ITeamInviteRepository interface {
	GetById(ctx context.Context, id string) (*models.TeamInvite, error)
	GetByToken(ctx context.Context, token string) (*models.TeamInvite, error)
	MarkAsAccepted(ctx context.Context, id string) error
	MarkAsCancelled(ctx context.Context, id string) error
	GetByEmailAndTeamId(ctx context.Context, email, teamId string) (*models.TeamInvite, error)
	Create(ctx context.Context, invite models.TeamInvite) (string, error)
	UpdateById(ctx context.Context, id string, invite models.TeamInvite) error
	DeleteById(ctx context.Context, id string) error
}

type teamInviteRepositoryParams struct {
	fx.In

	Database infrastructure.IDatabase
}

type teamInviteRepository struct {
	db *gorm.DB
}

func FxTeamInviteRepository() fx.Option {
	return fx_utils.AsProvider(newTeamInviteRepository, new(ITeamInviteRepository))
}

func newTeamInviteRepository(params teamInviteRepositoryParams) ITeamInviteRepository {
	return &teamInviteRepository{
		db: params.Database.GetInstance(),
	}
}

func (r *teamInviteRepository) GetById(ctx context.Context, id string) (*models.TeamInvite, error) {
	teamInvite := &models.TeamInvite{}
	err := r.db.WithContext(ctx).Preload("Team.Owner").Where("id = ?", id).First(teamInvite).Error

	if err != nil {
		if gorm_utils.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return teamInvite, nil
}

func (r *teamInviteRepository) GetByToken(ctx context.Context, token string) (*models.TeamInvite, error) {
	teamInvite := &models.TeamInvite{}
	err := r.db.WithContext(ctx).Where("token = ?", token).First(teamInvite).Error

	if err != nil {
		if gorm_utils.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return teamInvite, nil
}

func (r *teamInviteRepository) GetByEmailAndTeamId(ctx context.Context, email, teamId string) (*models.TeamInvite, error) {
	teamInvite := &models.TeamInvite{}
	err := r.db.WithContext(ctx).Where("team_id = ? AND email = ?", teamId, email).First(teamInvite).Error

	if err != nil {
		if gorm_utils.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return teamInvite, nil
}

func (r *teamInviteRepository) Create(ctx context.Context, invite models.TeamInvite) (string, error) {
	err := r.db.WithContext(ctx).Create(&invite).Error
	if err != nil {
		return "", err
	}

	return invite.Id, nil
}

func (r *teamInviteRepository) DeleteById(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.TeamInvite{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *teamInviteRepository) UpdateById(ctx context.Context, id string, invite models.TeamInvite) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Updates(&invite).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *teamInviteRepository) MarkAsAccepted(ctx context.Context, id string) error {
	return r.UpdateById(ctx, id, models.TeamInvite{
		Status: constants.TeamInviteStatusAccepted,
	})
}

func (r *teamInviteRepository) MarkAsCancelled(ctx context.Context, id string) error {
	return r.UpdateById(ctx, id, models.TeamInvite{
		Status: constants.TeamInviteStatusCancelled,
	})
}

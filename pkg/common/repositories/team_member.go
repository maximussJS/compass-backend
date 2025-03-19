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

type ITeamMemberRepository interface {
	Get(ctx context.Context, teamId, userId string) (*models.TeamMember, error)
	Create(ctx context.Context, member models.TeamMember) error
	GetAllByTeamId(ctx context.Context, teamId string) ([]models.TeamMember, error)
	Delete(ctx context.Context, teamId, userId string) error
}

type teamMemberRepositoryParams struct {
	fx.In

	Database infrastructure.IDatabase
}

type teamMemberRepository struct {
	db *gorm.DB
}

func FxTeamMemberRepository() fx.Option {
	return fx_utils.AsProvider(newTeamMemberRepository, new(ITeamMemberRepository))
}

func newTeamMemberRepository(params teamMemberRepositoryParams) ITeamMemberRepository {
	return &teamMemberRepository{
		db: params.Database.GetInstance(),
	}
}

func (r *teamMemberRepository) Get(ctx context.Context, teamId, userId string) (*models.TeamMember, error) {
	var teamMember *models.TeamMember
	err := r.db.WithContext(ctx).Where("team_id = ? AND user_id = ?", teamId, userId).First(teamMember).Error

	if err != nil {
		if gorm_utils.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("error getting team member: %w", err)
	}

	return teamMember, nil
}

func (r *teamMemberRepository) Create(ctx context.Context, member models.TeamMember) error {
	err := r.db.WithContext(ctx).Create(&member).Error
	if err != nil {
		return fmt.Errorf("failed to create team member: %w", err)
	}

	return nil
}

func (r *teamMemberRepository) GetAllByTeamId(ctx context.Context, teamId string) ([]models.TeamMember, error) {
	var teamMembers []models.TeamMember
	err := r.db.WithContext(ctx).Where("team_id = ?", teamId).Find(&teamMembers).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get team members by team id: %w", err)
	}

	return teamMembers, nil

}

func (r *teamMemberRepository) Delete(ctx context.Context, teamId, userId string) error {
	err := r.db.WithContext(ctx).Where("team_id = ? AND user_id = ?", teamId, userId).Delete(&models.TeamMember{}).Error

	if err != nil {
		return fmt.Errorf("failed to delete team member: %w", err)
	}

	return nil
}

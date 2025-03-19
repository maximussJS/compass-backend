package services

import (
	"compass-backend/pkg/api/api_errors"
	dto_team "compass-backend/pkg/api/common/dto/team"
	common_interfaces "compass-backend/pkg/api/common/interfaces"
	fx_utils "compass-backend/pkg/common/fx"
	common_lib "compass-backend/pkg/common/lib"
	common_models "compass-backend/pkg/common/models"
	common_repositories "compass-backend/pkg/common/repositories"
	"context"
	"go.uber.org/fx"
)

type ITeamService interface {
	common_interfaces.ITeamService
	Create(ctx context.Context, ownerId string, dto dto_team.CreateTeamRequest) (*common_models.Team, error)
	UpdateById(ctx context.Context, id string, dto dto_team.UpdateTeamRequest) (*common_models.Team, error)
	DeleteById(ctx context.Context, id string) error
}

type teamServiceParams struct {
	fx.In

	Logger         common_lib.ILogger
	TeamRepository common_repositories.ITeamRepository
}

type teamService struct {
	logger         common_lib.ILogger
	teamRepository common_repositories.ITeamRepository
}

func FxTeamService() fx.Option {
	return fx_utils.AsProvider(newTeamService, new(ITeamService))
}

func newTeamService(params teamServiceParams) ITeamService {
	return &teamService{
		logger:         params.Logger,
		teamRepository: params.TeamRepository,
	}
}

func (s *teamService) GetById(ctx context.Context, id string) (*common_models.Team, error) {
	team, err := s.teamRepository.GetById(ctx, id)

	if err != nil {
		s.logger.Errorf("failed to get team by id: %s", err)
		return nil, err
	}

	return team, nil
}

func (s *teamService) Create(ctx context.Context, ownerId string, dto dto_team.CreateTeamRequest) (*common_models.Team, error) {
	existingTeam, existingErr := s.teamRepository.GetByOwnerId(ctx, ownerId)

	if existingErr != nil {
		s.logger.Errorf("failed to get team by owner id %s %s", ownerId, existingErr)
		return nil, existingErr
	}

	if existingTeam != nil {
		return existingTeam, api_errors.ErrorOnlyOneTeamAllowed
	}

	id, createErr := s.teamRepository.Create(ctx, common_models.Team{
		OwnerId: ownerId,
		Name:    dto.Name,
	})

	if createErr != nil {
		s.logger.Errorf("failed to create team: %s", createErr)
		return nil, createErr
	}

	team, getErr := s.teamRepository.GetById(ctx, id)

	if getErr != nil {
		s.logger.Errorf("failed to get team by id: %s", getErr)
		return nil, getErr
	}

	return team, nil
}

func (s *teamService) UpdateById(ctx context.Context, id string, dto dto_team.UpdateTeamRequest) (*common_models.Team, error) {
	err := s.teamRepository.UpdateById(ctx, id, dto.ToModel())

	if err != nil {
		s.logger.Errorf("failed to update team: %s", err)
		return nil, err
	}

	team, err := s.teamRepository.GetById(ctx, id)

	if err != nil {
		s.logger.Errorf("failed to get team by id: %s", err)
		return nil, err
	}

	return team, nil
}

func (s *teamService) DeleteById(ctx context.Context, id string) error {
	err := s.teamRepository.DeleteById(ctx, id)

	if err != nil {
		s.logger.Errorf("failed to delete team: %s", err)
		return err
	}

	return nil
}

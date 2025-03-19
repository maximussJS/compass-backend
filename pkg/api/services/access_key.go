package services

import (
	"compass-backend/pkg/api/api_errors"
	dto_access_key "compass-backend/pkg/api/common/dto/access_key"
	fx_utils "compass-backend/pkg/common/fx"
	common_lib "compass-backend/pkg/common/lib"
	common_models "compass-backend/pkg/common/models"
	"compass-backend/pkg/common/repositories"
	"context"
	"go.uber.org/fx"
)

type IAccessKeyService interface {
	Create(ctx context.Context, userId string, team common_models.Team, dto dto_access_key.CreateAccessKeyRequest) (*common_models.AccessKey, error)
	GetById(ctx context.Context, userId, id string) (*common_models.AccessKey, error)
}

type accessKeyServiceParams struct {
	fx.In

	Logger              common_lib.ILogger
	TeamRepository      repositories.ITeamRepository
	AccessKeyRepository repositories.IAccessKeyRepository
}

type accessKeyService struct {
	logger    common_lib.ILogger
	team      repositories.ITeamRepository
	accessKey repositories.IAccessKeyRepository
}

func FxAccessKeyService() fx.Option {
	return fx_utils.AsProvider(newAccessKeyService, new(IAccessKeyService))
}

func newAccessKeyService(params accessKeyServiceParams) IAccessKeyService {
	return &accessKeyService{
		logger:    params.Logger,
		team:      params.TeamRepository,
		accessKey: params.AccessKeyRepository,
	}
}

func (s *accessKeyService) Create(ctx context.Context, userId string, team common_models.Team, dto dto_access_key.CreateAccessKeyRequest) (*common_models.AccessKey, error) {
	existingKey, err := s.accessKey.GetByNameAndTeamId(ctx, dto.Name, team.Id)

	if err != nil {
		s.logger.Errorf("failed to get access key by name and team id: %s", err)
		return nil, err
	}

	if existingKey != nil {
		return nil, api_errors.ErrorAccessKeyAlreadyExists
	}

	accessKey := common_models.AccessKey{
		UserId:    userId,
		TeamId:    team.Id,
		Name:      dto.Name,
		ExpiresAt: dto.ExpireTime(),
	}

	id, err := s.accessKey.Create(ctx, accessKey)
	if err != nil {
		s.logger.Errorf("failed to create access key: %s", err)
		return nil, err
	}

	newAccessKey, err := s.accessKey.GetById(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to get access key by id: %s", err)
		return nil, err
	}

	return newAccessKey, nil
}

func (s *accessKeyService) GetById(ctx context.Context, userId, id string) (*common_models.AccessKey, error) {
	accessKey, err := s.accessKey.GetByIdAndUserId(ctx, id, userId)

	if err != nil {
		s.logger.Errorf("failed to get access key by id: %s", err)
		return nil, err
	}

	if accessKey == nil {
		return nil, api_errors.ErrorAccessKeyNotFound
	}

	return accessKey, nil
}

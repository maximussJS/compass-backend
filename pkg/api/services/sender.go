package services

import (
	"compass-backend/pkg/api/lib"
	fx_utils "compass-backend/pkg/common/fx"
	common_infrastructure "compass-backend/pkg/common/infrastructure"
	common_types "compass-backend/pkg/common/types"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type ISenderService interface {
	SendTeamInvite(ctx context.Context, job common_types.SendTeamInviteEmailJob) error
	SendUserRegistered(ctx context.Context, job common_types.SendUserRegisteredEmailJob) error
}

type senderServiceParams struct {
	fx.In

	Env   lib.IEnv
	Redis common_infrastructure.IRedis
}

type senderService struct {
	userRegisteredRedisChannel string
	teamInviteRedisChannel     string
	redis                      *redis.Client
}

func FxSenderService() fx.Option {
	return fx_utils.AsProvider(newSenderService, new(ISenderService))
}

func newSenderService(params senderServiceParams) ISenderService {
	return &senderService{
		userRegisteredRedisChannel: params.Env.GetUserRegisteredRedisChannel(),
		teamInviteRedisChannel:     params.Env.GetTeamInviteRedisChannel(),
		redis:                      params.Redis.GetInstance(),
	}
}

func (s *senderService) SendTeamInvite(ctx context.Context, job common_types.SendTeamInviteEmailJob) error {
	payload, err := json.Marshal(job)

	if err != nil {
		return err
	}

	return s.publish(ctx, s.teamInviteRedisChannel, payload)
}

func (s *senderService) SendUserRegistered(ctx context.Context, job common_types.SendUserRegisteredEmailJob) error {
	payload, err := json.Marshal(job)

	if err != nil {
		return err
	}

	return s.publish(ctx, s.userRegisteredRedisChannel, payload)
}

func (s *senderService) publish(ctx context.Context, channel string, payload []byte) error {
	err := s.redis.Publish(ctx, channel, payload).Err()

	if err != nil {
		return err
	}

	return nil
}

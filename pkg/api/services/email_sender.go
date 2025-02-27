package services

import (
	"compass-backend/pkg/api/lib"
	"compass-backend/pkg/common/constants"
	fx_utils "compass-backend/pkg/common/fx"
	common_lib "compass-backend/pkg/common/pub_sub"
	common_types "compass-backend/pkg/common/types"
	"context"
	"go.uber.org/fx"
)

type IEmailSenderService interface {
	SendTeamInvite(ctx context.Context, data common_types.SendTeamInviteEmailJobData) error
	SendEmptyUserCreated(ctx context.Context, data common_types.SendEmptyUserCreatedEmailJobData) error
	SendUserRegistered(ctx context.Context, data common_types.SendUserRegisteredEmailJobData) error
}

type emailSenderServiceParams struct {
	fx.In

	Env            lib.IEnv
	RedisPublisher common_lib.IRedisPublisher
}

type emailSenderService struct {
	channel        string
	redisPublisher common_lib.IRedisPublisher
}

func FxEmailSenderService() fx.Option {
	return fx_utils.AsProvider(newEmailSenderService, new(IEmailSenderService))
}

func newEmailSenderService(params emailSenderServiceParams) IEmailSenderService {
	return &emailSenderService{
		channel:        params.Env.GetEmailRedisChannel(),
		redisPublisher: params.RedisPublisher,
	}
}

func (s *emailSenderService) SendTeamInvite(ctx context.Context, data common_types.SendTeamInviteEmailJobData) error {
	job := common_types.EmailJob{
		Type: constants.SendTeamInviteEmailJobType,
		Data: data,
	}

	return s.redisPublisher.Publish(ctx, s.channel, job)
}

func (s *emailSenderService) SendEmptyUserCreated(ctx context.Context, data common_types.SendEmptyUserCreatedEmailJobData) error {
	job := common_types.EmailJob{
		Type: constants.SendEmptyUserCreatedEmailJobType,
		Data: data,
	}

	return s.redisPublisher.Publish(ctx, s.channel, job)
}

func (s *emailSenderService) SendUserRegistered(ctx context.Context, data common_types.SendUserRegisteredEmailJobData) error {
	job := common_types.EmailJob{
		Type: constants.SendUserRegisteredEmailJobType,
		Data: data,
	}

	return s.redisPublisher.Publish(ctx, s.channel, job)
}

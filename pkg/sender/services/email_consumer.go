package services

import (
	"compass-backend/pkg/common/constants"
	fx_utils "compass-backend/pkg/common/fx"
	common_lib "compass-backend/pkg/common/lib"
	"compass-backend/pkg/common/pub_sub"
	common_types "compass-backend/pkg/common/types"
	"compass-backend/pkg/sender/lib"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type IEmailConsumerService interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type emailConsumerServiceParams struct {
	fx.In

	Logger            common_lib.ILogger
	Env               lib.IEnv
	TeamInviteService ITeamInviteService
	UserService       IUserService
	RedisConsumer     pub_sub.IRedisConsumer
}

type emailConsumerService struct {
	channel           string
	logger            common_lib.ILogger
	redisConsumer     pub_sub.IRedisConsumer
	teamInviteService ITeamInviteService
	userService       IUserService
}

func FxEmailConsumerService() fx.Option {
	return fx_utils.AsProvider(newEmailConsumerService, new(IEmailConsumerService))
}

func newEmailConsumerService(lc fx.Lifecycle, params emailConsumerServiceParams) IEmailConsumerService {
	emailConsumer := &emailConsumerService{
		logger:            params.Logger,
		channel:           params.Env.GetEmailRedisChannel(),
		redisConsumer:     params.RedisConsumer,
		teamInviteService: params.TeamInviteService,
		userService:       params.UserService,
	}

	lc.Append(fx.Hook{
		OnStart: emailConsumer.Start,
		OnStop:  emailConsumer.Stop,
	})

	return emailConsumer
}

func (s *emailConsumerService) Start(ctx context.Context) error {
	return s.redisConsumer.Start(ctx, s.channel, s.consumerFn())
}

func (s *emailConsumerService) Stop(ctx context.Context) error {
	return s.redisConsumer.Stop(ctx)
}

func (s *emailConsumerService) consumerFn() pub_sub.RedisConsumerFn {
	return func(ctx context.Context, msg *redis.Message) error {
		var job common_types.EmailJob

		if err := json.Unmarshal([]byte(msg.Payload), &job); err != nil {
			s.logger.Error(fmt.Sprintf("error unmarshalling team invite message: %v", err))
			return err
		}

		dataBytes, err := json.Marshal(job.Data)
		if err != nil {
			s.logger.Error(fmt.Sprintf("error marshalling job.Data: %v", err))
			return err
		}

		switch job.Type {
		case constants.SendTeamInviteEmailJobType:
			{
				var data common_types.SendTeamInviteEmailJobData

				if err := json.Unmarshal(dataBytes, &data); err != nil {
					s.logger.Error(fmt.Sprintf("error unmarshalling team invite message: %v", err))
					return err
				}

				if err := s.teamInviteService.SendTeamInvite(ctx, data.Id, data.AcceptLink, data.CancelLink); err != nil {
					s.logger.Error(fmt.Sprintf("error sending team invite: %v", err))
					return err
				}

				s.logger.Info(fmt.Sprintf("Team invite email sent: %v", data.Id))

				return nil
			}
		case constants.SendEmptyUserCreatedEmailJobType:
			{
				var data common_types.SendEmptyUserCreatedEmailJobData

				if err := json.Unmarshal(dataBytes, &data); err != nil {
					s.logger.Error(fmt.Sprintf("error unmarshalling empty user created message: %v", err))
					return err
				}

				if err := s.userService.SendEmptyUserCreated(ctx, data.Email, data.Password); err != nil {
					s.logger.Error(fmt.Sprintf("error sending empty user created email: %v", err))
					return err
				}

				s.logger.Info(fmt.Sprintf("empty user created email sent to: %v", data.Email))

				return nil
			}
		case constants.SendConfirmEmailJobType:
			{
				var data common_types.SendConfirmEmailJobData

				if err := json.Unmarshal(dataBytes, &data); err != nil {
					s.logger.Error(fmt.Sprintf("error unmarshalling empty user created message: %v", err))
					return err
				}

				if err := s.userService.SendConfirmEmail(ctx, data.Email, data.Name, data.ConfirmationLink); err != nil {
					s.logger.Error(fmt.Sprintf("error sending empty user created email: %v", err))
					return err
				}

				s.logger.Info(fmt.Sprintf("empty user created email sent to: %v", data.Email))

				return nil
			}
		default:
			{
				s.logger.Error(fmt.Sprintf("unknown job type: %v", job.Type))
				return nil
			}
		}
	}
}

package processors

import (
	"compass-backend/pkg/sender/interfaces"
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	fx_utils "compass-backend/pkg/common/fx"
	"compass-backend/pkg/common/infrastructure"
	common_lib "compass-backend/pkg/common/lib"
	"compass-backend/pkg/common/types"
	"compass-backend/pkg/sender/lib"
	"compass-backend/pkg/sender/services"
)

type IUserRegisteredProcessor interface {
	interfaces.IJobProcessor
}

type userRegisteredProcessorParams struct {
	fx.In

	Env          lib.IEnv
	Logger       common_lib.ILogger
	Redis        infrastructure.IRedis
	UserRegister services.IUserRegisterService
}

type userRegisteredProcessor struct {
	base         *BaseProcessor
	userRegister services.IUserRegisterService
}

func FxUserRegisteredProcessor() fx.Option {
	return fx_utils.AsProvider(newUserRegisteredProcessor, new(IUserRegisteredProcessor))
}

func newUserRegisteredProcessor(lc fx.Lifecycle, params userRegisteredProcessorParams) IUserRegisteredProcessor {
	processor := &userRegisteredProcessor{
		userRegister: params.UserRegister,
	}
	base := NewBaseProcessor(
		"UserRegisteredProcessor",
		params.Redis.GetInstance(),
		params.Env.GetUserRegisteredChannel(),
		params.Logger,
		processor.Process,
	)
	processor.base = base

	lc.Append(fx.Hook{
		OnStart: base.Start,
		OnStop:  base.Stop,
	})

	return processor
}

func (p *userRegisteredProcessor) Process(ctx context.Context, msg *redis.Message) error {
	var job types.SendEmptyUserCreatedEmailJob

	if err := json.Unmarshal([]byte(msg.Payload), &job); err != nil {
		p.base.logger.Error(fmt.Sprintf("error unmarshalling user registered message: %v", err))
		return err
	}

	if err := p.userRegister.SendRegisterUser(ctx, job.Email, job.Password); err != nil {
		p.base.logger.Error(fmt.Sprintf("error sending user registration email: %v", err))
		return err
	}

	p.base.logger.Info(fmt.Sprintf("User registration email sent to: %v", job.Email))

	return nil
}

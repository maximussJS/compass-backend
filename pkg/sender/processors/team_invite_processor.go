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

type ITeamInviteProcessor interface {
	interfaces.IJobProcessor
}

type teamInviteProcessorParams struct {
	fx.In

	Env               lib.IEnv
	Logger            common_lib.ILogger
	Redis             infrastructure.IRedis
	TeamInviteService services.ITeamInviteService
}

type teamInviteProcessor struct {
	base              *BaseProcessor
	teamInviteService services.ITeamInviteService
}

func FxTeamInviteProcessor() fx.Option {
	return fx_utils.AsProvider(newTeamInviteProcessor, new(ITeamInviteProcessor))
}

func newTeamInviteProcessor(lc fx.Lifecycle, params teamInviteProcessorParams) ITeamInviteProcessor {
	processor := &teamInviteProcessor{
		teamInviteService: params.TeamInviteService,
	}

	base := NewBaseProcessor(
		"TeamInviteProcessor",
		params.Redis.GetInstance(),
		params.Env.GetTeamInviteRedisChannel(),
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

func (p *teamInviteProcessor) Process(ctx context.Context, msg *redis.Message) error {
	var job types.SendTeamInviteEmailJob

	if err := json.Unmarshal([]byte(msg.Payload), &job); err != nil {
		p.base.logger.Error(fmt.Sprintf("error unmarshalling team invite message: %v", err))
		return err
	}

	if err := p.teamInviteService.SendTeamInvite(ctx, job.Id, job.AcceptLink, job.CancelLink); err != nil {
		p.base.logger.Error(fmt.Sprintf("error sending team invite: %v", err))
		return err
	}

	p.base.logger.Info(fmt.Sprintf("Team invite email sent: %v", job.Id))

	return nil
}

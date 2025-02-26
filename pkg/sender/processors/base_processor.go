package processors

import (
	"context"
	"fmt"

	common_lib "compass-backend/pkg/common/lib"
	"github.com/redis/go-redis/v9"
)

type ProcessFunc func(ctx context.Context, msg *redis.Message) error

type BaseProcessor struct {
	name             string
	cancelProcessing context.CancelFunc
	redis            *redis.Client
	pubSub           *redis.PubSub
	channelName      string
	logger           common_lib.ILogger
	processFn        ProcessFunc
}

func NewBaseProcessor(
	name string,
	redis *redis.Client,
	channelName string,
	logger common_lib.ILogger,
	processFn ProcessFunc,
) *BaseProcessor {
	return &BaseProcessor{
		name:        name,
		redis:       redis,
		channelName: channelName,
		logger:      logger,
		processFn:   processFn,
	}
}

func (bp *BaseProcessor) Start(ctx context.Context) error {
	bp.pubSub = bp.redis.Subscribe(ctx, bp.channelName)
	processCtx, cancel := context.WithCancel(context.Background())
	bp.cancelProcessing = cancel

	go bp.startProcessing(processCtx)
	bp.logger.Info(fmt.Sprintf("%s started on channel %s", bp.name, bp.channelName))
	return nil
}

func (bp *BaseProcessor) startProcessing(ctx context.Context) {
	ch := bp.pubSub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			err := bp.processFn(ctx, msg)

			if err != nil {
				bp.logger.Error(fmt.Sprintf("%s error processing message: %v", bp.name, err))
			}
		}
	}
}

func (bp *BaseProcessor) Stop(_ context.Context) error {
	if bp.cancelProcessing != nil {
		bp.cancelProcessing()
	}
	if err := bp.pubSub.Close(); err != nil {
		return fmt.Errorf("%s stop error: %w", bp.name, err)
	}
	bp.logger.Info(fmt.Sprintf("%s stopped on channel %s", bp.name, bp.channelName))

	return nil
}

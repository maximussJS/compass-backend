package pub_sub

import (
	fx_utils "compass-backend/pkg/common/fx"
	common_infrastructure "compass-backend/pkg/common/infrastructure"
	"compass-backend/pkg/common/lib"
	"context"
	"fmt"
	"go.uber.org/fx"

	"github.com/redis/go-redis/v9"
)

type RedisConsumerFn func(ctx context.Context, msg *redis.Message) error

type IRedisConsumer interface {
	Start(ctx context.Context, channel string, fn RedisConsumerFn) error
	Stop(ctx context.Context) error
}

type redisConsumerParams struct {
	fx.In

	Env    lib.IEnv
	Logger lib.ILogger
	Redis  common_infrastructure.IRedis
}

type redisConsumer struct {
	redis      *redis.Client
	pubSub     *redis.PubSub
	logger     lib.ILogger
	cancelFunc context.CancelFunc
}

func FxRedisConsumer() fx.Option {
	return fx_utils.AsProvider(newRedisConsumer, new(IRedisConsumer))
}

func newRedisConsumer(params redisConsumerParams) IRedisConsumer {
	return &redisConsumer{
		redis:  params.Redis.GetInstance(),
		logger: params.Logger,
	}
}

func (c *redisConsumer) Start(ctx context.Context, channel string, fn RedisConsumerFn) error {
	c.pubSub = c.redis.Subscribe(ctx, channel)

	processCtx, cancel := context.WithCancel(context.Background())

	c.cancelFunc = cancel

	go c.startProcessing(processCtx, fn)

	c.logger.Info(fmt.Sprintf("redis consumer started on channel %s", channel))
	return nil
}

func (c *redisConsumer) startProcessing(ctx context.Context, fn RedisConsumerFn) {
	ch := c.pubSub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			err := fn(ctx, msg)

			if err != nil {
				c.logger.Error(fmt.Sprintf("redis consumer error processing message: %v", err))
			}
		}
	}
}

func (c *redisConsumer) Stop(_ context.Context) error {
	if err := c.pubSub.Close(); err != nil {
		return fmt.Errorf("redis consumer stop error: %w", err)
	}

	c.logger.Info(fmt.Sprint("redis consumer stopped"))

	return nil
}

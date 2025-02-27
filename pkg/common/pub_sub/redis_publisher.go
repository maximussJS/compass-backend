package pub_sub

import (
	fx_utils "compass-backend/pkg/common/fx"
	common_infrastructure "compass-backend/pkg/common/infrastructure"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type IRedisPublisher interface {
	Publish(ctx context.Context, channel string, data interface{}) error
}

type redisPublisherParams struct {
	fx.In

	Redis common_infrastructure.IRedis
}

type redisPublisher struct {
	redis *redis.Client
}

func FxRedisPublisher() fx.Option {
	return fx_utils.AsProvider(newRedisPublisher, new(IRedisPublisher))
}

func newRedisPublisher(params redisPublisherParams) IRedisPublisher {
	return &redisPublisher{
		redis: params.Redis.GetInstance(),
	}
}

func (s *redisPublisher) Publish(ctx context.Context, channel string, data interface{}) error {
	payload, err := json.Marshal(data)

	if err != nil {
		return err
	}

	err = s.redis.Publish(ctx, channel, payload).Err()

	if err != nil {
		return err
	}

	return nil
}

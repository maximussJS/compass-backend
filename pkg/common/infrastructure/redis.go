package infrastructure

import (
	fx_utils "compass-backend/pkg/common/fx"
	"compass-backend/pkg/common/lib"
	"context"
	"fmt"
	go_redis "github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type IRedis interface {
	GetInstance() *go_redis.Client
}

type redisParams struct {
	fx.In

	Env    lib.IEnv
	Logger lib.ILogger
}

type redis struct {
	env    lib.IEnv
	logger lib.ILogger
	client *go_redis.Client
}

func FxRedis() fx.Option {
	return fx_utils.AsProvider(newRedis, new(IRedis))
}

func newRedis(lc fx.Lifecycle, params redisParams) IRedis {
	redis := &redis{
		env:    params.Env,
		logger: params.Logger,
	}
	opt, err := go_redis.ParseURL(params.Env.GetRedisUrl())

	if err != nil {
		redis.logger.Error(fmt.Sprintf("Failed to parse redis url: %s", err))
		panic(err)
	}

	client := go_redis.NewClient(opt)

	ctx := context.TODO()

	_, err = client.Ping(ctx).Result()

	if err != nil {
		redis.logger.Error(fmt.Sprintf("Failed to connect to redis: %s", err))
		panic(err)
	}

	redis.logger.Info("Connected to redis")

	redis.client = client

	lc.Append(fx.Hook{
		OnStop: redis.shutdown,
	})

	return redis
}

func (r *redis) GetInstance() *go_redis.Client {
	return r.client
}

func (r *redis) shutdown(_ context.Context) error {
	closeErr := r.client.Close()

	if closeErr != nil {
		return fmt.Errorf("failed to close redis connection: %s", closeErr)
	}

	r.logger.Info("Redis connection closed")

	return nil
}

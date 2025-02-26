package interfaces

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type IJobProcessor interface {
	Process(ctx context.Context, msg *redis.Message) error
}

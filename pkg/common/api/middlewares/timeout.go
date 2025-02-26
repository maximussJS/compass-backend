package middlewares

import (
	fx_utils "compass-backend/pkg/common/fx"
	"compass-backend/pkg/common/lib"
	"context"
	"go.uber.org/fx"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ITimeoutMiddleware interface {
	Handle() gin.HandlerFunc
}

type timeoutMiddlewareParams struct {
	fx.In

	Env lib.IEnv
}

type timeoutMiddleware struct {
	timeoutDuration time.Duration
}

func FxTimeoutMiddleware() fx.Option {
	return fx_utils.AsProvider(newTimeoutMiddleware, new(ITimeoutMiddleware))
}

func newTimeoutMiddleware(params timeoutMiddlewareParams) ITimeoutMiddleware {
	return &timeoutMiddleware{
		timeoutDuration: params.Env.GetRequestTimeoutDuration(),
	}
}

func (p timeoutMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), p.timeoutDuration)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		finished := make(chan struct{}, 1)

		go func() {
			c.Next()
			finished <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{
				"error": "request timeout",
			})
		case <-finished:
		}
	}
}

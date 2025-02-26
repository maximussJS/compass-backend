package middlewares

import (
	"compass-backend/pkg/api/api/responses"
	gin_utils "compass-backend/pkg/api/utils/gin"
	fx_utils "compass-backend/pkg/common/fx"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ITrainerMiddleware interface {
	Handle() gin.HandlerFunc
}

type trainerMiddleware struct {
}

func FxTrainerMiddleware() fx.Option {
	return fx_utils.AsProvider(newTrainerMiddleware, new(ITrainerMiddleware))
}

func newTrainerMiddleware() ITrainerMiddleware {
	return &trainerMiddleware{}
}

func (m trainerMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := gin_utils.GetUser(c)

		if err != nil {
			responses.UnauthorizedWithAbort(c)
			return
		}

		if !user.IsTrainer() {
			responses.ForbiddenWithAbort(c)
			return
		}

		c.Next()
	}
}

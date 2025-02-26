package middlewares

import (
	"compass-backend/pkg/api/api/responses"
	gin_utils "compass-backend/pkg/api/utils/gin"
	fx_utils "compass-backend/pkg/common/fx"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type IClientMiddleware interface {
	Handle() gin.HandlerFunc
}

type clientMiddleware struct {
}

func FxClientMiddleware() fx.Option {
	return fx_utils.AsProvider(newClientMiddleware, new(IClientMiddleware))
}

func newClientMiddleware() IClientMiddleware {
	return &clientMiddleware{}
}

func (m clientMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := gin_utils.GetUser(c)

		if err != nil {
			responses.UnauthorizedWithAbort(c)
			return
		}

		if !user.IsClient() {
			responses.ForbiddenWithAbort(c)
			return
		}

		c.Next()
	}
}

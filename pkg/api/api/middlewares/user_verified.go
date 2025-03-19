package middlewares

import (
	"compass-backend/pkg/api/api/responses"
	gin_utils "compass-backend/pkg/api/utils/gin"
	fx_utils "compass-backend/pkg/common/fx"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type IUserVerifiedMiddleware interface {
	Handle() gin.HandlerFunc
}

type userVerifiedMiddleware struct {
}

func FxUserVerifiedMiddleware() fx.Option {
	return fx_utils.AsProvider(newUserVerifiedMiddleware, new(IUserVerifiedMiddleware))
}

func newUserVerifiedMiddleware() IUserVerifiedMiddleware {
	return &userVerifiedMiddleware{}
}

func (m userVerifiedMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := gin_utils.GetUser(c)

		if err != nil {
			responses.UnauthorizedWithAbort(c)
			return
		}

		if !user.IsVerified {
			responses.ForbiddenWithMessage(c, "User is not verified")
			return
		}

		c.Next()
	}
}

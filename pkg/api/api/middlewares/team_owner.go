package middlewares

import (
	"compass-backend/pkg/api/api/responses"
	gin_utils "compass-backend/pkg/api/utils/gin"
	fx_utils "compass-backend/pkg/common/fx"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ITeamOwnerMiddleware interface {
	Handle() gin.HandlerFunc
}

type teamOwnerMiddleware struct{}

func FxTeamOwnerMiddleware() fx.Option {
	return fx_utils.AsProvider(newTeamOwnerMiddleware, new(ITeamOwnerMiddleware))
}

func newTeamOwnerMiddleware() ITeamOwnerMiddleware {
	return &teamOwnerMiddleware{}
}

func (m teamOwnerMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := gin_utils.GetUser(c)

		if err != nil {
			responses.UnauthorizedWithAbort(c)
			return
		}

		team, err := gin_utils.GetTeam(c)

		if err != nil {
			responses.NotFoundWithAbort(c, "Team not found")
			return
		}

		if !team.IsOwner(user) {
			responses.ForbiddenWithAbort(c)
			return
		}

		c.Next()
	}
}

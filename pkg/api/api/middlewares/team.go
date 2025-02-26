package middlewares

import (
	"compass-backend/pkg/api/api/responses"
	common_interfaces "compass-backend/pkg/api/common/interfaces"
	gin_utils "compass-backend/pkg/api/utils/gin"
	fx_utils "compass-backend/pkg/common/fx"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ITeamMiddleware interface {
	Handle() gin.HandlerFunc
	SetTeamService(teamService common_interfaces.ITeamService)
}

type teamMiddleware struct {
	teamService common_interfaces.ITeamService
}

func FxTeamMiddleware() fx.Option {
	return fx_utils.AsProvider(newTeamMiddleware, new(ITeamMiddleware))
}

func newTeamMiddleware() ITeamMiddleware {
	return &teamMiddleware{}
}

func (m *teamMiddleware) SetTeamService(teamService common_interfaces.ITeamService) {
	m.teamService = teamService
}

func (m teamMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		teamId, ok := gin_utils.GetStringTeamIdParam(c)
		if !ok {
			return
		}

		team, err := m.teamService.GetById(c, teamId)

		if err != nil {
			responses.InternalServerError(c)
			return
		}

		if team == nil {
			responses.NotFoundWithAbort(c, fmt.Sprintf("Team with id %s not found", teamId))
			return
		}

		gin_utils.SetTeam(c, team)

		c.Next()
	}
}

package controllers

import (
	"compass-backend/pkg/api/api/responses"
	"compass-backend/pkg/api/api_errors"
	dto_team "compass-backend/pkg/api/common/dto/team"
	"compass-backend/pkg/api/services"
	gin_utils "compass-backend/pkg/api/utils/gin"
	fx_utils "compass-backend/pkg/common/fx"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ITeamController interface {
	Create(c *gin.Context)
	GetById(c *gin.Context)
	UpdateById(c *gin.Context)
	DeleteById(c *gin.Context)
}

type teamControllerParams struct {
	fx.In

	TeamService services.ITeamService
}

type teamController struct {
	teamService services.ITeamService
}

func FxTeamController() fx.Option {
	return fx_utils.AsProvider(newTeamController, new(ITeamController))
}

func newTeamController(params teamControllerParams) ITeamController {
	return &teamController{
		teamService: params.TeamService,
	}
}

func (h *teamController) Create(c *gin.Context) {
	user, userErr := gin_utils.GetUser(c)

	if userErr != nil {
		responses.UnauthorizedWithAbort(c)
		return
	}

	var req dto_team.CreateTeamRequest

	if ok := gin_utils.BindData(c, &req); !ok {
		return
	}

	team, err := h.teamService.Create(c, user.Id, req)

	if err != nil {
		if errors.Is(err, api_errors.ErrorOnlyOneTeamAllowed) {
			responses.BadRequestWithAbort(c, "Only one team is allowed")
			return
		}
		responses.InternalServerError(c)
		return
	}

	responses.SuccessJsonWithMessage(c, team, fmt.Sprintf("Team %s created", team.Name))
}

func (h *teamController) GetById(c *gin.Context) {
	team, err := gin_utils.GetTeam(c)

	if err != nil {
		responses.NotFoundWithAbort(c, "Team not found")
		return
	}

	responses.SuccessJson(c, team)
}

func (h *teamController) UpdateById(c *gin.Context) {
	team, err := gin_utils.GetTeam(c)

	if err != nil {
		responses.NotFoundWithAbort(c, "Team not found")
		return
	}

	var req dto_team.UpdateTeamRequest

	if ok := gin_utils.BindData(c, &req); !ok {
		return
	}

	updatedTeam, err := h.teamService.UpdateById(c, team.Id, req)

	if err != nil {
		responses.InternalServerError(c)
		return
	}

	responses.SuccessJsonWithMessage(c, updatedTeam, fmt.Sprintf("Team %s updated", updatedTeam.Name))
}

func (h *teamController) DeleteById(c *gin.Context) {
	team, err := gin_utils.GetTeam(c)

	if err != nil {
		responses.NotFoundWithAbort(c, "Team not found")
		return
	}

	err = h.teamService.DeleteById(c, team.Id)

	if err != nil {
		responses.InternalServerError(c)
		return
	}

	responses.SuccessJsonWithMessage(c, team, fmt.Sprintf("Team %s deleted", team.Name))
}

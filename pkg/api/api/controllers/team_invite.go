package controllers

import (
	"compass-backend/pkg/api/api/responses"
	"compass-backend/pkg/api/api_errors"
	invite_dto "compass-backend/pkg/api/common/dto/invite"
	"compass-backend/pkg/api/services"
	gin_utils "compass-backend/pkg/api/utils/gin"
	fx_utils "compass-backend/pkg/common/fx"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ITeamInviteController interface {
	InviteByEmail(c *gin.Context)
	AcceptInvite(c *gin.Context)
	CancelInvite(c *gin.Context)
}

type teamInviteControllerParams struct {
	fx.In

	InviteService services.ITeamInviteService
}

type teamInviteController struct {
	inviteService services.ITeamInviteService
}

func FxTeamInviteController() fx.Option {
	return fx_utils.AsProvider(newTeamInviteController, new(ITeamInviteController))
}

func newTeamInviteController(params teamInviteControllerParams) ITeamInviteController {
	return &teamInviteController{
		inviteService: params.InviteService,
	}
}

func (h *teamInviteController) InviteByEmail(c *gin.Context) {
	user, userErr := gin_utils.GetUser(c)

	if userErr != nil {
		responses.UnauthorizedWithAbort(c)
		return
	}

	var req invite_dto.InviteByEmailRequest

	if ok := gin_utils.BindData(c, &req); !ok {
		return
	}

	err := h.inviteService.InviteByEmail(c, req.Email, user.Id)

	if err != nil {
		if errors.Is(err, api_errors.ErrorTeamInviteAlreadySend) {
			responses.BadRequestWithAbort(c, fmt.Sprintf("Invitation to %s already sent", req.Email))
			return
		}

		if errors.Is(err, api_errors.ErrorTeamNotFound) {
			responses.NotFoundWithAbort(c, fmt.Sprintf("You doesn't have any team yet"))
			return
		}

		responses.InternalServerError(c)
		return
	}

	responses.SuccessMessage(c, fmt.Sprintf("Invitation to %s has been sent", req.Email))
}

func (h *teamInviteController) AcceptInvite(c *gin.Context) {
	token, ok := gin_utils.GetTokenParam(c)
	if !ok {
		return
	}

	err := h.inviteService.AcceptInvite(c, token)

	if err != nil {
		if errors.Is(err, api_errors.ErrorInvalidToken) {
			responses.NotFoundWithAbort(c, "Invalid token")
			return
		}
		responses.InternalServerError(c)
		return
	}

	responses.SuccessMessage(c, "Invitation has been accepted")
}

func (h *teamInviteController) CancelInvite(c *gin.Context) {
	token, ok := gin_utils.GetTokenParam(c)
	if !ok {
		return
	}

	err := h.inviteService.CancelInvite(c, token)

	if err != nil {
		if errors.Is(err, api_errors.ErrorInvalidToken) {
			responses.NotFoundWithAbort(c, "Invalid token")
			return
		}
		responses.InternalServerError(c)
		return
	}

	responses.SuccessMessage(c, "Invitation has been canceled")
}

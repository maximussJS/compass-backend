package controllers

import (
	"compass-backend/pkg/api/api/responses"
	"compass-backend/pkg/api/api_errors"
	dto_access_key "compass-backend/pkg/api/common/dto/access_key"
	"compass-backend/pkg/api/lib"
	"compass-backend/pkg/api/services"
	gin_utils "compass-backend/pkg/api/utils/gin"
	fx_utils "compass-backend/pkg/common/fx"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type IAccessKeyController interface {
	Create(c *gin.Context)
	GetById(c *gin.Context)
	ListForTeam(c *gin.Context)
	UpdateById(c *gin.Context)
	DeleteById(c *gin.Context)
}

type accessKeyControllerParams struct {
	fx.In

	Env              lib.IEnv
	AccessKeyService services.IAccessKeyService
}

type accessKeyController struct {
	defaultLimit int
	maxLimit     int
	service      services.IAccessKeyService
}

func FxAccessKeyController() fx.Option {
	return fx_utils.AsProvider(newAccessKeyController, new(IAccessKeyController))
}

func newAccessKeyController(params accessKeyControllerParams) IAccessKeyController {
	return &accessKeyController{
		service:      params.AccessKeyService,
		defaultLimit: params.Env.GetDefaultLimit(),
		maxLimit:     params.Env.GetMaxLimit(),
	}
}

func (h *accessKeyController) Create(c *gin.Context) {
	user, userErr := gin_utils.GetUser(c)

	if userErr != nil {
		responses.UnauthorizedWithAbort(c)
		return
	}

	team, teamErr := gin_utils.GetTeam(c)
	if teamErr != nil {
		responses.NotFoundWithAbort(c, "Team not found")
		return
	}

	var req dto_access_key.CreateAccessKeyRequest

	if ok := gin_utils.BindData(c, &req); !ok {
		return
	}

	key, err := h.service.Create(c, user.Id, *team, req)

	if err != nil {
		if errors.Is(err, api_errors.ErrorTeamNotFound) {
			responses.BadRequestWithAbort(c, fmt.Sprintf("Team with id %s not found", team.Id))
			return
		}

		if errors.Is(err, api_errors.ErrorAccessKeyAlreadyExists) {
			responses.BadRequestWithAbort(c, fmt.Sprintf("Access key with name %s already exists", req.Name))
			return
		}

		responses.InternalServerError(c)
		return
	}

	responses.SuccessJsonWithMessage(c, key, fmt.Sprintf("Access Key %s created", key.Name))
}

func (h *accessKeyController) GetById(c *gin.Context) {
	user, userErr := gin_utils.GetUser(c)

	if userErr != nil {
		responses.UnauthorizedWithAbort(c)
		return
	}

	id, ok := gin_utils.GetStringIdParam(c)
	if !ok {
		responses.BadRequest(c, fmt.Sprintf("Invalid id %s", id))
		return
	}

	key, err := h.service.GetById(c, user.Id, id)

	if err != nil {
		if errors.Is(err, api_errors.ErrorAccessKeyNotFound) {
			responses.NotFound(c, fmt.Sprintf("Access key with id %s not found", id))
			return
		}

		responses.InternalServerError(c)
		return
	}

	responses.SuccessJson(c, key)
}

func (h *accessKeyController) ListForTeam(c *gin.Context) {
	user, userErr := gin_utils.GetUser(c)

	if userErr != nil {
		responses.UnauthorizedWithAbort(c)
		return
	}

	team, teamErr := gin_utils.GetTeam(c)
	if teamErr != nil {
		responses.NotFoundWithAbort(c, "Team not found")
		return
	}

	limit := gin_utils.GetLimit(c, h.defaultLimit, h.maxLimit)
	offset := gin_utils.GetOffset(c)

	keys, err := h.service.ListForTeam(c, user.Id, team, limit, offset)

	if err != nil {
		responses.InternalServerError(c)
		return
	}

	responses.SuccessJson(c, keys)
}

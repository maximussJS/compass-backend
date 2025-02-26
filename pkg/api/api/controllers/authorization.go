package controllers

import (
	"compass-backend/pkg/api/api/responses"
	"compass-backend/pkg/api/api_errors"
	authorization_dto "compass-backend/pkg/api/common/dto/authorization"
	"compass-backend/pkg/api/services"
	gin_utils "compass-backend/pkg/api/utils/gin"
	fx_utils "compass-backend/pkg/common/fx"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type IAuthorizationController interface {
	SignInByPassword(c *gin.Context)
	SignUpByPassword(c *gin.Context)
}

type authorizationControllerParams struct {
	fx.In

	AuthorizationService services.IAuthorizationService
}

type authorizationController struct {
	authorizationService services.IAuthorizationService
}

func FxAuthorizationController() fx.Option {
	return fx_utils.AsProvider(newAuthorizationController, new(IAuthorizationController))
}

func newAuthorizationController(params authorizationControllerParams) IAuthorizationController {
	return &authorizationController{
		authorizationService: params.AuthorizationService,
	}
}

func (h *authorizationController) SignInByPassword(c *gin.Context) {
	var req authorization_dto.SignInByPasswordRequest

	if ok := gin_utils.BindData(c, &req); !ok {
		return
	}

	token, err := h.authorizationService.SignInByPassword(c, req)

	if err != nil {
		if errors.Is(err, api_errors.ErrorInvalidCredentials) {
			responses.BadRequest(c, "Invalid credentials")
			return
		}

		responses.InternalServerError(c)
		return
	}

	responses.SuccessJson(c, token)
}

func (h *authorizationController) SignUpByPassword(c *gin.Context) {
	var req authorization_dto.SignUpByPasswordRequest

	if ok := gin_utils.BindData(c, &req); !ok {
		return
	}

	user, err := h.authorizationService.SignUpByPassword(c, req)

	if err != nil {
		if errors.Is(err, api_errors.ErrorUserAlreadyExists) {
			responses.BadRequest(c, fmt.Sprintf("User with email %s already exists", req.Email))
			return
		}

		responses.InternalServerError(c)
		return
	}

	responses.SuccessJsonWithMessage(c, user, fmt.Sprintf("User with email %s successfully created", req.Email))
}

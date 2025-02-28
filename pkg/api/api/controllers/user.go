package controllers

import (
	"compass-backend/pkg/api/api/responses"
	"compass-backend/pkg/api/api_errors"
	user_dto "compass-backend/pkg/api/common/dto/user"
	"compass-backend/pkg/api/services"
	gin_utils "compass-backend/pkg/api/utils/gin"
	fx_utils "compass-backend/pkg/common/fx"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type IUserController interface {
	ChangeName(c *gin.Context)
	ChangePassword(c *gin.Context)
	ConfirmEmail(c *gin.Context)
	Me(c *gin.Context)
}

type userControllerParams struct {
	fx.In

	UserService services.IUserService
}

type userController struct {
	userService services.IUserService
}

func FxUserController() fx.Option {
	return fx_utils.AsProvider(newUserController, new(IUserController))
}

func newUserController(params userControllerParams) IUserController {
	return &userController{
		userService: params.UserService,
	}
}

func (h *userController) ChangeName(c *gin.Context) {
	user, userErr := gin_utils.GetUser(c)

	if userErr != nil {
		responses.UnauthorizedWithAbort(c)
		return
	}

	var req user_dto.ChangeNameRequest

	if ok := gin_utils.BindData(c, &req); !ok {
		return
	}

	err := h.userService.ChangeName(c, user, req)

	if err != nil {
		if errors.Is(err, api_errors.ErrorNameIsTheSame) {
			responses.BadRequest(c, "New name is the same as the old one")
			return
		}

		responses.InternalServerError(c)
		return
	}

	responses.SuccessMessage(c, fmt.Sprintf("Name changed to %s", req.Name))
}

func (h *userController) ChangePassword(c *gin.Context) {
	user, userErr := gin_utils.GetUser(c)

	if userErr != nil {
		responses.UnauthorizedWithAbort(c)
		return
	}

	var req user_dto.ChangePasswordRequest

	if ok := gin_utils.BindData(c, &req); !ok {
		return
	}

	err := h.userService.ChangePassword(c, user, req)

	if err != nil {
		if errors.Is(err, api_errors.ErrorInvalidCredentials) {
			responses.Unauthorized(c)
			return
		}

		if errors.Is(err, api_errors.ErrorUserPasswordIsTheSame) {
			responses.BadRequest(c, "New password is the same as the old one")
			return
		}

		responses.InternalServerError(c)
		return
	}

	responses.SuccessMessage(c, "Password changed")
}

func (h *userController) Me(c *gin.Context) {
	user, userErr := gin_utils.GetUser(c)

	if userErr != nil {
		responses.UnauthorizedWithAbort(c)
		return
	}

	responses.SuccessJson(c, user)
}

func (h *userController) ConfirmEmail(c *gin.Context) {
	token, ok := gin_utils.GetTokenParam(c)
	if !ok {
		return
	}

	err := h.userService.ConfirmEmail(c, token)

	if err != nil {
		if errors.Is(err, api_errors.ErrorEmailAlreadyConfirmed) {
			responses.BadRequest(c, "Email already confirmed")
			return
		}

		if errors.Is(err, api_errors.ErrorInvalidToken) {
			responses.BadRequest(c, "Invalid token")
			return
		}

		if errors.Is(err, api_errors.ErrorUserNotFound) {
			responses.NotFound(c, "User not found")
			return
		}

		responses.InternalServerError(c)
		return
	}

	responses.SuccessMessage(c, "Email confirmed")
}

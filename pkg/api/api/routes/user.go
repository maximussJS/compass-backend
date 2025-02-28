package routes

import (
	"compass-backend/pkg/api/api/controllers"
	"compass-backend/pkg/api/api/middlewares"
	common_routes "compass-backend/pkg/common/api/routes"
	common_infrastracture "compass-backend/pkg/common/infrastructure"
	common_lib "compass-backend/pkg/common/lib"
	"fmt"
	"go.uber.org/fx"
)

type userRoute struct {
	logger                  common_lib.ILogger
	router                  common_infrastracture.IRouter
	userController          controllers.IUserController
	authorizationMiddleware middlewares.IAuthorizationMiddleware
}

type userRouteParams struct {
	fx.In

	Logger                  common_lib.ILogger
	Router                  common_infrastracture.IRouter
	UserController          controllers.IUserController
	AuthorizationMiddleware middlewares.IAuthorizationMiddleware
}

func FxUserRoute() fx.Option {
	return common_routes.AsRoute(newUserRoute)
}

func newUserRoute(params userRouteParams) common_routes.IRoute {
	return &userRoute{
		logger:                  params.Logger,
		router:                  params.Router,
		userController:          params.UserController,
		authorizationMiddleware: params.AuthorizationMiddleware,
	}
}

func (h *userRoute) Setup() {
	group := h.router.GetRouter().Group("/api/users")

	h.logger.Info(fmt.Sprintf("Mapped User Route %s", group.BasePath()))

	group.GET("/confirm-email/:token", h.userController.ConfirmEmail)

	h.logger.Info(fmt.Sprintf("Mapped GET %s/confirm-email/:token", group.BasePath()))

	group.GET("/me", h.authorizationMiddleware.Handle(), h.userController.Me)

	h.logger.Info(fmt.Sprintf("Mapped GET %s/me", group.BasePath()))

	group.PATCH("/change-name", h.authorizationMiddleware.Handle(), h.userController.ChangeName)

	h.logger.Info(fmt.Sprintf("Mapped POST %s/change-name", group.BasePath()))

	group.PATCH("/change-password", h.authorizationMiddleware.Handle(), h.userController.ChangePassword)

	h.logger.Info(fmt.Sprintf("Mapped POST %s/change-password", group.BasePath()))
}

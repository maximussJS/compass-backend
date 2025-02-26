package routes

import (
	"compass-backend/pkg/api/api/controllers"
	common_routes "compass-backend/pkg/common/api/routes"
	common_infrastracture "compass-backend/pkg/common/infrastructure"
	common_lib "compass-backend/pkg/common/lib"
	"fmt"
	"go.uber.org/fx"
)

type authorizationRoute struct {
	logger                  common_lib.ILogger
	router                  common_infrastracture.IRouter
	authorizationController controllers.IAuthorizationController
}

type authorizationRouteParams struct {
	fx.In

	Logger                  common_lib.ILogger
	Router                  common_infrastracture.IRouter
	AuthorizationController controllers.IAuthorizationController
}

func FxAuthorizationRoute() fx.Option {
	return common_routes.AsRoute(newAuthorizationRoute)
}

func newAuthorizationRoute(params authorizationRouteParams) common_routes.IRoute {
	return &authorizationRoute{
		logger:                  params.Logger,
		router:                  params.Router,
		authorizationController: params.AuthorizationController,
	}
}

func (h *authorizationRoute) Setup() {
	group := h.router.GetRouter().Group("api/authorization")

	h.logger.Info(fmt.Sprintf("Mapped Authorization Route %s", group.BasePath()))

	group.POST("/sign-in-by-password", h.authorizationController.SignInByPassword)
	h.logger.Info(fmt.Sprintf("Mapped POST %s/sign-in-by-password ", group.BasePath()))
	group.POST("/sign-up-by-password", h.authorizationController.SignUpByPassword)
	h.logger.Info(fmt.Sprintf("Mapped POST %s/sign-up-by-password ", group.BasePath()))
}

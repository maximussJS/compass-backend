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

type teamInviteRoute struct {
	logger                  common_lib.ILogger
	router                  common_infrastracture.IRouter
	teamInviteController    controllers.ITeamInviteController
	authorizationMiddleware middlewares.IAuthorizationMiddleware
	trainerMiddleware       middlewares.ITrainerMiddleware
}

type teamInviteRouteParams struct {
	fx.In

	Logger                  common_lib.ILogger
	TeamInviteController    controllers.ITeamInviteController
	AuthorizationMiddleware middlewares.IAuthorizationMiddleware
	TrainerMiddleware       middlewares.ITrainerMiddleware
	Router                  common_infrastracture.IRouter
}

func FxTeamInviteRoute() fx.Option {
	return common_routes.AsRoute(newTeamInviteRoute)
}

func newTeamInviteRoute(params teamInviteRouteParams) common_routes.IRoute {
	return &teamInviteRoute{
		logger:                  params.Logger,
		router:                  params.Router,
		teamInviteController:    params.TeamInviteController,
		authorizationMiddleware: params.AuthorizationMiddleware,
		trainerMiddleware:       params.TrainerMiddleware,
	}
}

func (h *teamInviteRoute) Setup() {
	group := h.router.GetRouter().Group("/api/team-invites")

	h.logger.Info(fmt.Sprintf("Mapped Invite Route %s", group.BasePath()))

	group.POST("/send-by-email",
		h.authorizationMiddleware.Handle(),
		h.trainerMiddleware.Handle(),
		h.teamInviteController.InviteByEmail,
	)

	h.logger.Info(fmt.Sprintf("Mapped POST %s/send-by-email", group.BasePath()))

	group.GET("/accept/:token", h.teamInviteController.AcceptInvite)
	h.logger.Info(fmt.Sprintf("Mapped GET %s/accept/:token", group.BasePath()))

	group.GET("/cancel/:token", h.teamInviteController.CancelInvite)
	h.logger.Info(fmt.Sprintf("Mapped GET %s/cancel/:token", group.BasePath()))
}

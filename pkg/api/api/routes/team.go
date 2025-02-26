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

type teamRoute struct {
	logger                  common_lib.ILogger
	router                  common_infrastracture.IRouter
	teamController          controllers.ITeamController
	authorizationMiddleware middlewares.IAuthorizationMiddleware
	trainerMiddleware       middlewares.ITrainerMiddleware
	teamMiddleware          middlewares.ITeamMiddleware
	teamOwnerMiddleware     middlewares.ITeamOwnerMiddleware
}

type teamRouteParams struct {
	fx.In

	Logger                  common_lib.ILogger
	Router                  common_infrastracture.IRouter
	TeamController          controllers.ITeamController
	AuthorizationMiddleware middlewares.IAuthorizationMiddleware
	TrainerMiddleware       middlewares.ITrainerMiddleware
	TeamMiddleware          middlewares.ITeamMiddleware
	TeamOwnerMiddleware     middlewares.ITeamOwnerMiddleware
}

func FxTeamRoute() fx.Option {
	return common_routes.AsRoute(newTeamRoute)
}

func newTeamRoute(params teamRouteParams) common_routes.IRoute {
	return &teamRoute{
		logger:                  params.Logger,
		router:                  params.Router,
		teamController:          params.TeamController,
		authorizationMiddleware: params.AuthorizationMiddleware,
		trainerMiddleware:       params.TrainerMiddleware,
		teamMiddleware:          params.TeamMiddleware,
		teamOwnerMiddleware:     params.TeamOwnerMiddleware,
	}
}

func (h *teamRoute) Setup() {
	group := h.router.GetRouter().Group("/api/teams")

	group.Use(h.authorizationMiddleware.Handle())
	group.Use(h.trainerMiddleware.Handle())

	h.logger.Info(fmt.Sprintf("Mapped Team Route %s", group.BasePath()))

	group.POST("", h.teamController.Create)

	h.logger.Info(fmt.Sprintf("Mapped POST %s", group.BasePath()))

	group.GET(
		"/:teamId",
		h.teamMiddleware.Handle(),
		h.teamOwnerMiddleware.Handle(),
		h.teamController.GetById,
	)

	h.logger.Info(fmt.Sprintf("Mapped GET %s/:teamId", group.BasePath()))

	group.PATCH(
		"/:teamId",
		h.teamMiddleware.Handle(),
		h.teamOwnerMiddleware.Handle(),
		h.teamController.UpdateById,
	)

	h.logger.Info(fmt.Sprintf("Mapped PATCH %s/:teamId", group.BasePath()))

	group.DELETE(
		"/:teamId",
		h.teamMiddleware.Handle(),
		h.teamOwnerMiddleware.Handle(),
		h.teamController.DeleteById,
	)

	h.logger.Info(fmt.Sprintf("Mapped DELETE %s/:teamId", group.BasePath()))
}

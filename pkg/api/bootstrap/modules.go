package bootstrap

import (
	"compass-backend/pkg/api/api/controllers"
	"compass-backend/pkg/api/api/middlewares"
	"compass-backend/pkg/api/api/routes"
	"compass-backend/pkg/api/infrastructure"
	"compass-backend/pkg/api/lib"
	"compass-backend/pkg/api/repositories"
	"compass-backend/pkg/api/services"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	lib.Module,
	infrastructure.Module,
	controllers.Module,
	routes.Module,
	services.Module,
	repositories.Module,
	middlewares.Module,
	fx.Invoke(func(
		authorizationMiddleware middlewares.IAuthorizationMiddleware,
		teamMiddleware middlewares.ITeamMiddleware,
		authorizationService services.IAuthorizationService,
		teamService services.ITeamService,
	) {
		authorizationMiddleware.SetTokenService(authorizationService)
		teamMiddleware.SetTeamService(teamService)
	}),
)

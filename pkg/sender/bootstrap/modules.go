package bootstrap

import (
	"compass-backend/pkg/sender/api/controllers"
	"compass-backend/pkg/sender/api/middlewares"
	"compass-backend/pkg/sender/api/routes"
	"compass-backend/pkg/sender/infrastructure"
	"compass-backend/pkg/sender/lib"
	"compass-backend/pkg/sender/repositories"
	"compass-backend/pkg/sender/services"
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
)

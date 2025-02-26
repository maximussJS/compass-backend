package routes

import (
	common_routes "compass-backend/pkg/common/api/routes"
	"go.uber.org/fx"
)

var Module = fx.Options(
	common_routes.FxHealthCheckRoute(),
	common_routes.FxRoutes(),
	fx.Invoke(func(routes common_routes.Routes) {}),
)

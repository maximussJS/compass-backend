package routes

import (
	"compass-backend/pkg/common/api/controllers"
	"compass-backend/pkg/common/infrastructure"
	"compass-backend/pkg/common/lib"
	"fmt"
	"go.uber.org/fx"
)

type healthCheckRoute struct {
	logger                lib.ILogger
	router                infrastructure.IRouter
	healthCheckController controllers.IHealthCheckController
}

type healthCheckRouteParams struct {
	fx.In

	Logger                lib.ILogger
	HealthCheckController controllers.IHealthCheckController
	Router                infrastructure.IRouter
}

func FxHealthCheckRoute() fx.Option {
	return AsRoute(newHealthCheckRoute)
}

func newHealthCheckRoute(params healthCheckRouteParams) IRoute {
	return &healthCheckRoute{
		logger:                params.Logger,
		router:                params.Router,
		healthCheckController: params.HealthCheckController,
	}
}

func (h *healthCheckRoute) Setup() {
	group := h.router.GetRouter().Group("/health-check")

	h.logger.Info(fmt.Sprintf("Mapped HealthCheck Route %s", group.BasePath()))

	group.GET("", h.healthCheckController.HealthCheck)

	h.logger.Info(fmt.Sprintf("Mapped GET %s", group.BasePath()))
}

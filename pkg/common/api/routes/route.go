package routes

import (
	fx_utils "compass-backend/pkg/common/fx"
	"context"
	"go.uber.org/fx"
)

type routesParams struct {
	fx.In

	Routes []IRoute `group:"routes"`
}

type Routes []IRoute

type IRoute interface {
	Setup()
}

func FxRoutes() fx.Option {
	return fx.Provide(newRoutes)
}

func AsRoute(r any) fx.Option {
	return fx_utils.AsNamedProvider(r, new(IRoute), `group:"routes"`)
}

func newRoutes(lc fx.Lifecycle, params routesParams) Routes {
	routes := Routes(params.Routes)

	lc.Append(fx.Hook{
		OnStart: routes.Setup,
	})

	return routes
}

func (r Routes) Setup(context.Context) error {
	for _, route := range r {
		route.Setup()
	}

	return nil
}

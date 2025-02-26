package bootstrap

import (
	"compass-backend/pkg/common/infrastructure"
	"go.uber.org/fx"
)

func CreateApp() fx.Option {
	return fx.Options(
		Modules,
		fx.Invoke(func(server infrastructure.IHttpServer) {}),
	)
}

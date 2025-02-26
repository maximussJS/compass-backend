package infrastructure

import (
	common_infrastracture "compass-backend/pkg/common/infrastructure"
	"go.uber.org/fx"
)

var Module = fx.Options(
	common_infrastracture.FxRouter(),
	common_infrastracture.FxRedis(),
	common_infrastracture.FxDatabase(),
	common_infrastracture.FxHttpServer(),
)

package controllers

import (
	common_controllers "compass-backend/pkg/common/api/controllers"
	"go.uber.org/fx"
)

var Module = fx.Options(
	common_controllers.FxHealthCheckController(),
)

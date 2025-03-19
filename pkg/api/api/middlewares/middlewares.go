package middlewares

import (
	"compass-backend/pkg/common/api/middlewares"
	"go.uber.org/fx"
)

var Module = fx.Options(
	middlewares.FxTimeoutMiddleware(),
	FxAuthorizationMiddleware(),
	FxTeamMiddleware(),
	FxTeamOwnerMiddleware(),
	FxUserVerifiedMiddleware(),
)

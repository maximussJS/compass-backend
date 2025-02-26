package repositories

import (
	"compass-backend/pkg/common/repositories"
	"go.uber.org/fx"
)

var Module = fx.Options(
	repositories.FxTeamInviteRepository(),
)

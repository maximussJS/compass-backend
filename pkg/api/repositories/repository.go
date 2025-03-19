package repositories

import (
	"compass-backend/pkg/common/repositories"
	"go.uber.org/fx"
)

var Module = fx.Options(
	repositories.FxUserRepository(),
	repositories.FxTeamRepository(),
	repositories.FxTeamInviteRepository(),
	repositories.FxTeamMemberRepository(),
	repositories.FxAccessKeyRepository(),
)

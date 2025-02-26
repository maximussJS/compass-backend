package repositories

import (
	"compass-backend/pkg/common/repositories"
	"go.uber.org/fx"
)

var Module = fx.Options(
	FxCategoryRepository(),
	repositories.FxUserRepository(),
	repositories.FxTeamRepository(),
	repositories.FxTeamInviteRepository(),
	repositories.FxTeamMemberRepository(),
)

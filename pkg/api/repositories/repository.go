package repositories

import (
	"compass-backend/pkg/common/repositories"
	"go.uber.org/fx"
)

var Module = fx.Options(
	FxCategoryRepository(),
	FxExerciseRepository(),
	FxExerciseMediaRepository(),
	repositories.FxUserRepository(),
	repositories.FxTeamRepository(),
	repositories.FxTeamInviteRepository(),
	repositories.FxTeamMemberRepository(),
)

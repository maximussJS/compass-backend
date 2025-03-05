package services

import (
	"compass-backend/pkg/common/pub_sub"
	"go.uber.org/fx"
)

var Module = fx.Options(
	pub_sub.FxRedisPublisher(),
	FxEmailSenderService(),
	FxUserService(),
	FxTokenService(),
	FxAuthorizationService(),
	FxCategoryService(),
	FxTeamInviteService(),
	FxTeamService(),
	FxExerciseService(),
)

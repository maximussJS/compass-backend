package services

import (
	"compass-backend/pkg/common/services"
	"go.uber.org/fx"
)

var Module = fx.Options(
	FxAuthorizationService(),
	FxCategoryService(),
	FxTeamInviteService(),
	FxTeamService(),
	services.FxEmailSenderService(),
	FxUserService(),
)

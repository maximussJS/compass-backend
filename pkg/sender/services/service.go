package services

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	FxMailService(),
	FxTeamInviteService(),
	FxUserRegisterService(),
)

package services

import (
	"compass-backend/pkg/common/pub_sub"
	"go.uber.org/fx"
)

var Module = fx.Options(
	pub_sub.FxRedisConsumer(),
	FxMailService(),
	FxTeamInviteService(),
	FxUserService(),
	FxEmailConsumerService(),
	fx.Invoke(func(service IEmailConsumerService) {}),
)

package processors

import "go.uber.org/fx"

var Module = fx.Options(
	FxTeamInviteProcessor(),
	FxUserRegisteredProcessor(),
	fx.Invoke(func(
		teamInviteProcessor ITeamInviteProcessor,
		userRegisteredProcessor IUserRegisteredProcessor,
	) {}),
)

package api_errors

import "errors"

var (
	ErrorTeamInviteAlreadySend = errors.New("team invite already send")
	ErrorTeamInviteExpired     = errors.New("team invite expired")
	ErrorTeamInviteAccepted    = errors.New("team invite already accepted")
	ErrorTeamInviteCancelled   = errors.New("team invite cancelled")
	ErrorTeamInviteNotFound    = errors.New("team invite not found")
)

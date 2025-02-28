package constants

type EmailJobType string

const (
	SendTeamInviteEmailJobType       EmailJobType = "send_team_invite_email"
	SendEmptyUserCreatedEmailJobType EmailJobType = "send_empty_user_created_email"
	SendConfirmEmailJobType          EmailJobType = "send_confirm_email"
)

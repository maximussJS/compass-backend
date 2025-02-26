package constants

type TeamInviteStatus string

const (
	TeamInviteStatusPending   TeamInviteStatus = "pending"
	TeamInviteStatusAccepted  TeamInviteStatus = "accepted"
	TeamInviteStatusCancelled TeamInviteStatus = "cancelled"
)

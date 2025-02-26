package types

type SendTeamInviteEmailJob struct {
	Id         string `json:"id"`
	AcceptLink string `json:"acceptLink"`
	CancelLink string `json:"cancelLink"`
}

type SendUserRegisteredEmailJob struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

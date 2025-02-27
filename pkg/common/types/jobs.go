package types

import "compass-backend/pkg/common/constants"

type EmailJob struct {
	Type constants.EmailJobType
	Data interface{}
}

type SendTeamInviteEmailJobData struct {
	Id         string `json:"id"`
	AcceptLink string `json:"acceptLink"`
	CancelLink string `json:"cancelLink"`
}

type SendEmptyUserCreatedEmailJobData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SendUserRegisteredEmailJobData struct {
	Name string `json:"name"`
}

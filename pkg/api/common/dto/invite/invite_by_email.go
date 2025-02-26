package invite

type InviteByEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

package user

type ChangePasswordRequest struct {
	Password    string `json:"password" binding:"required,gte=6,lte=30"`
	OldPassword string `json:"oldPassword" binding:"required,gte=6,lte=30"`
}

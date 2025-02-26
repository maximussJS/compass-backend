package authorization

type SignUpByPasswordRequest struct {
	Name     string `json:"name" binding:"required,gte=1,lte=30"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

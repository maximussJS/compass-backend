package user

type ChangeNameRequest struct {
	Name string `json:"name" binding:"required,gte=1,lte=30"`
}

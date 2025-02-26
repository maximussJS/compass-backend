package team

type CreateTeamRequest struct {
	Name string `json:"name" binding:"required,gte=1,lte=30"`
}

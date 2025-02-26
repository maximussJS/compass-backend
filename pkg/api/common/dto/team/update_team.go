package team

import "compass-backend/pkg/common/models"

type UpdateTeamRequest struct {
	Name string `json:"name" binding:"required,gte=1,lte=30"`
}

func (r *UpdateTeamRequest) ToModel() models.Team {
	return models.Team{
		Name: r.Name,
	}
}

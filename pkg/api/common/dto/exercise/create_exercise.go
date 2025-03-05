package exercise

import "compass-backend/pkg/api/models"

type CreateExerciseRequest struct {
	Name        string `json:"name" binding:"required,gte=3,lte=30"`
	Description string `json:"description" binding:"required,gte=3,lte=30"`
}

func (c *CreateExerciseRequest) ToModel() models.Exercise {
	return models.Exercise{
		Name:        c.Name,
		Description: c.Description,
	}
}

package category

import "compass-backend/pkg/api/models"

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,gte=3,lte=30"`
}

func (c *CreateCategoryRequest) ToModel() models.Category {
	return models.Category{
		Name: c.Name,
	}
}

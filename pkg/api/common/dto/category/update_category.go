package category

import "compass-backend/pkg/api/models"

type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"gte=3,lte=30"`
}

func (c *UpdateCategoryRequest) ToModel() models.Category {
	return models.Category{
		Name: c.Name,
	}
}

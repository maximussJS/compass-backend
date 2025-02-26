package interfaces

import (
	"compass-backend/pkg/common/models"
	"context"
)

type ITeamService interface {
	GetById(ctx context.Context, id string) (*models.Team, error)
}

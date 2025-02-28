package interfaces

import (
	"compass-backend/pkg/common/models"
	"context"
)

type IAuthService interface {
	GetUserByToken(ctx context.Context, token string) (*models.User, error)
}

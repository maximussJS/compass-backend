package claims

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type ConfirmEmailClaims struct {
	UserId string `json:"userId"`
	*jwt.RegisteredClaims
}

func (c *ConfirmEmailClaims) IsExpired() bool {
	return c.ExpiresAt != nil && c.ExpiresAt.Before(time.Now())
}

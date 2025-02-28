package claims

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AuthClaims struct {
	UserId string `json:"user_id"`
	*jwt.RegisteredClaims
}

func (c AuthClaims) String() string {
	return fmt.Sprintf("AuthClaims{UserId: %s, RegisteredClaims: %v}", c.UserId, c.RegisteredClaims)
}

func (c *AuthClaims) IsExpired() bool {
	return c.ExpiresAt != nil && c.ExpiresAt.Before(time.Now())
}

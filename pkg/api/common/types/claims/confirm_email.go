package claims

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type InviteClaims struct {
	Email  string `json:"email"`
	TeamId string `json:"team_id"`
	*jwt.RegisteredClaims
}

func (c InviteClaims) String() string {
	return fmt.Sprintf("InviteClaims{Email: %s, TeamId: %s, RegisteredClaims: %v}", c.Email, c.TeamId, c.RegisteredClaims)
}

func (c *InviteClaims) IsExpired() bool {
	return c.ExpiresAt != nil && c.ExpiresAt.Before(time.Now())
}

package claims

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	UserId string `json:"user_id"`
	*jwt.RegisteredClaims
}

func (c AuthClaims) String() string {
	return fmt.Sprintf("AuthClaims{UserId: %s, RegisteredClaims: %v}", c.UserId, c.RegisteredClaims)
}

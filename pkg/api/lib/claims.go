package lib

import (
	common_types "compass-backend/pkg/api/common/types/claims"
	"compass-backend/pkg/common/constants"
	fx_utils "compass-backend/pkg/common/fx"
	go_jwt "github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
	"time"
)

type IClaims interface {
	NewAuthClaims(userId string, role constants.UserRole) common_types.AuthClaims
	NewInviteClaims(email, trainerId string) common_types.InviteClaims
}

type claimsParams struct {
	fx.In

	Env IEnv
}

type claims struct {
	issuer                   string
	authExpirationDuration   time.Duration
	inviteExpirationDuration time.Duration
}

func FxClaims() fx.Option {
	return fx_utils.AsProvider(newClaims, new(IClaims))
}

func newClaims(params claimsParams) *claims {
	return &claims{
		issuer:                   params.Env.GetAppName(),
		authExpirationDuration:   params.Env.GetAuthExpirationDuration(),
		inviteExpirationDuration: params.Env.GetInviteExpirationDuration(),
	}
}

func (c *claims) NewAuthClaims(userId string, role constants.UserRole) common_types.AuthClaims {
	return common_types.AuthClaims{
		UserId: userId,
		RegisteredClaims: &go_jwt.RegisteredClaims{
			Subject:   userId,
			Issuer:    c.issuer,
			Audience:  []string{string(role)},
			ExpiresAt: go_jwt.NewNumericDate(time.Now().Add(c.authExpirationDuration)),
			IssuedAt:  go_jwt.NewNumericDate(time.Now()),
		},
	}
}

func (c *claims) NewInviteClaims(email, teamId string) common_types.InviteClaims {
	return common_types.InviteClaims{
		Email:  email,
		TeamId: teamId,
		RegisteredClaims: &go_jwt.RegisteredClaims{
			Subject:   email,
			Issuer:    c.issuer,
			ExpiresAt: go_jwt.NewNumericDate(time.Now().Add(c.inviteExpirationDuration)),
			IssuedAt:  go_jwt.NewNumericDate(time.Now()),
		},
	}
}

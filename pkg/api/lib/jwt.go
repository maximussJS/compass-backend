package lib

import (
	fx_utils "compass-backend/pkg/common/fx"
	common_lib "compass-backend/pkg/common/lib"
	"fmt"
	go_jwt "github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
)

type IJwt interface {
	Generate(claims go_jwt.Claims) (string, error)
	Verify(tokenString string, customClaims go_jwt.Claims) error
}

type jwtParams struct {
	fx.In

	Logger common_lib.ILogger
	Env    IEnv
}

type jwt struct {
	logger        common_lib.ILogger
	secret        []byte
	signingMethod *go_jwt.SigningMethodHMAC
}

func FxJwt() fx.Option {
	return fx_utils.AsProvider(newJwt, new(IJwt))
}

func newJwt(params jwtParams) *jwt {
	return &jwt{
		logger:        params.Logger,
		secret:        params.Env.GetJwtSecret(),
		signingMethod: go_jwt.SigningMethodHS512,
	}
}

func (s *jwt) Generate(claims go_jwt.Claims) (string, error) {
	token := go_jwt.NewWithClaims(s.signingMethod, claims)

	tokenString, err := token.SignedString(s.secret)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to generate token: %s", err))
		return "", err
	}

	return tokenString, nil
}

func (s *jwt) Verify(tokenString string, customClaims go_jwt.Claims) error {
	token, err := go_jwt.ParseWithClaims(tokenString, customClaims, func(token *go_jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*go_jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

package middlewares

import (
	"compass-backend/pkg/api/api/responses"
	"compass-backend/pkg/api/api_errors"
	common_interfaces "compass-backend/pkg/api/common/interfaces"
	gin_utils "compass-backend/pkg/api/utils/gin"
	fx_utils "compass-backend/pkg/common/fx"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"strings"
)

type IAuthorizationMiddleware interface {
	Handle() gin.HandlerFunc
	SetTokenService(tokenService common_interfaces.ITokenService)
}

type authorizationMiddlewareParams struct {
	fx.In
}

type authorizationMiddleware struct {
	tokenService common_interfaces.ITokenService
}

func FxAuthorizationMiddleware() fx.Option {
	return fx_utils.AsProvider(newAuthorizationMiddleware, new(IAuthorizationMiddleware))
}

func newAuthorizationMiddleware(_ authorizationMiddlewareParams) IAuthorizationMiddleware {
	return &authorizationMiddleware{}
}

func (m authorizationMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := m.getTokenFromRequest(c)

		if err != nil {
			responses.UnauthorizedWithAbort(c)
			return
		}

		user, err := m.tokenService.GetUserByToken(c, token)

		if err != nil {
			if errors.Is(err, api_errors.ErrorInvalidToken) {
				responses.UnauthorizedWithAbort(c)
				return
			}

			responses.InternalServerError(c)
			return
		}

		gin_utils.SetUser(c, user)

		c.Next()
	}
}

func (m *authorizationMiddleware) getTokenFromRequest(c *gin.Context) (string, error) {
	tokenString := c.Request.Header.Get("Authorization")

	if tokenString == "" {
		return "", api_errors.ErrorUnauthorized
	}

	tokenParts := strings.Split(tokenString, " ")

	if len(tokenParts) != 2 {
		return "", api_errors.ErrorInvalidToken
	}

	return tokenParts[1], nil
}

func (m *authorizationMiddleware) SetTokenService(tokenService common_interfaces.ITokenService) {
	m.tokenService = tokenService
}

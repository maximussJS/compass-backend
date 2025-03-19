package services

import (
	"compass-backend/pkg/api/api_errors"
	"compass-backend/pkg/api/common/types/claims"
	"compass-backend/pkg/api/lib"
	fx_utils "compass-backend/pkg/common/fx"
	common_lib "compass-backend/pkg/common/lib"
	"compass-backend/pkg/common/models"
	"go.uber.org/fx"
)

type ITokenService interface {
	GenerateAuthorizationToken(user *models.User) (string, error)
	VerifyAuthorizationToken(token string) (*claims.AuthClaims, error)
	GenerateTeamInviteToken(teamId, email string) (token string, expiresAt int64, err error)
	VerifyTeamInviteToken(token string) (*claims.InviteClaims, error)
	GenerateConfirmEmailToken(userId string) (string, error)
	VerifyConfirmEmailToken(token string) (*claims.ConfirmEmailClaims, error)
}

type tokenServiceParams struct {
	fx.In

	Logger common_lib.ILogger
	Jwt    lib.IJwt
	Claims lib.IClaims
}

type tokenService struct {
	logger common_lib.ILogger
	jwt    lib.IJwt
	claims lib.IClaims
}

func FxTokenService() fx.Option {
	return fx_utils.AsProvider(newTokenService, new(ITokenService))
}

func newTokenService(params tokenServiceParams) ITokenService {
	return &tokenService{
		logger: params.Logger,
		jwt:    params.Jwt,
		claims: params.Claims,
	}
}

func (s *tokenService) GenerateAuthorizationToken(user *models.User) (string, error) {
	authClaims := s.claims.NewAuthClaims(user.Id, user.Role)

	token, tokenErr := s.jwt.Generate(authClaims)

	if tokenErr != nil {
		s.logger.Errorf("failed to generate auth token: %s", tokenErr)
		return "", tokenErr
	}

	return token, nil
}

func (s *tokenService) VerifyAuthorizationToken(token string) (*claims.AuthClaims, error) {
	var authClaims claims.AuthClaims

	err := s.jwt.Verify(token, &authClaims)

	if err != nil || authClaims.UserId == "" {
		return nil, api_errors.ErrorInvalidToken
	}

	return &authClaims, nil
}

func (s *tokenService) GenerateTeamInviteToken(teamId, email string) (token string, expiresAt int64, err error) {
	inviteClaims := s.claims.NewInviteClaims(email, teamId)

	token, err = s.jwt.Generate(inviteClaims)

	if err != nil {
		s.logger.Errorf("failed to generate invite token: %s", err)
		return "", 0, err
	}

	expiresAt = inviteClaims.ExpiresAt.Unix()

	return
}

func (s *tokenService) VerifyTeamInviteToken(token string) (*claims.InviteClaims, error) {
	var inviteClaims claims.InviteClaims

	err := s.jwt.Verify(token, &inviteClaims)

	if err != nil || inviteClaims.Email == "" {
		return nil, api_errors.ErrorInvalidToken
	}

	return &inviteClaims, nil
}

func (s *tokenService) GenerateConfirmEmailToken(userId string) (string, error) {
	confirmEmailClaims := s.claims.NewConfirmEmailClaims(userId)

	token, tokenErr := s.jwt.Generate(confirmEmailClaims)

	if tokenErr != nil {
		s.logger.Errorf("failed to generate confirm email token: %s", tokenErr)
		return "", tokenErr
	}

	return token, nil
}

func (s *tokenService) VerifyConfirmEmailToken(token string) (*claims.ConfirmEmailClaims, error) {
	var confirmEmailClaims claims.ConfirmEmailClaims

	err := s.jwt.Verify(token, &confirmEmailClaims)

	if err != nil || confirmEmailClaims.UserId == "" {
		return nil, api_errors.ErrorInvalidToken
	}

	return &confirmEmailClaims, nil
}

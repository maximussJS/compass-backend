package services

import (
	"compass-backend/pkg/api/api_errors"
	authorization_dto "compass-backend/pkg/api/common/dto/authorization"
	common_interfaces "compass-backend/pkg/api/common/interfaces"
	common_types "compass-backend/pkg/api/common/types/claims"
	"compass-backend/pkg/api/lib"
	crypto_utils "compass-backend/pkg/api/utils/crypto"
	fx_utils "compass-backend/pkg/common/fx"
	common_lib "compass-backend/pkg/common/lib"
	common_models "compass-backend/pkg/common/models"
	"context"
	"fmt"
	"go.uber.org/fx"
)

type IAuthorizationService interface {
	common_interfaces.ITokenService
	SignInByPassword(ctx context.Context, dto authorization_dto.SignInByPasswordRequest) (string, error)
	SignUpByPassword(ctx context.Context, dto authorization_dto.SignUpByPasswordRequest) (*common_models.User, error)
}

type authorizationServiceParams struct {
	fx.In

	Logger      common_lib.ILogger
	Jwt         lib.IJwt
	Claims      lib.IClaims
	userService IUserService
}

type authorizationService struct {
	logger      common_lib.ILogger
	jwt         lib.IJwt
	claims      lib.IClaims
	userService IUserService
}

func FxAuthorizationService() fx.Option {
	return fx_utils.AsProvider(newAuthorizationService, new(IAuthorizationService))
}

func newAuthorizationService(params authorizationServiceParams) *authorizationService {
	return &authorizationService{
		logger:      params.Logger,
		jwt:         params.Jwt,
		claims:      params.Claims,
		userService: params.userService,
	}
}

func (s *authorizationService) SignInByPassword(ctx context.Context, dto authorization_dto.SignInByPasswordRequest) (string, error) {
	user, userErr := s.userService.GetByEmail(ctx, dto.Email)

	if userErr != nil {
		s.logger.Error(fmt.Sprintf("failed to get user by email: %s", userErr))
		return "", userErr
	}

	if user == nil {
		return "", api_errors.ErrorInvalidCredentials
	}

	passwordIsValid := crypto_utils.VerifyHash(dto.Password, user.Password)

	if !passwordIsValid {
		return "", api_errors.ErrorInvalidCredentials
	}

	authClaims := s.claims.NewAuthClaims(user.Id, user.Role)

	token, tokenErr := s.jwt.Generate(authClaims)

	if tokenErr != nil {
		s.logger.Error(fmt.Sprintf("failed to generate auth token: %s", tokenErr))
		return "", tokenErr
	}

	return token, nil
}

func (s *authorizationService) SignUpByPassword(ctx context.Context, dto authorization_dto.SignUpByPasswordRequest) (*common_models.User, error) {
	existingUser, userErr := s.userService.GetByEmail(ctx, dto.Email)

	if userErr != nil {
		s.logger.Error(fmt.Sprintf("failed to get user by email: %s", userErr))
		return nil, userErr
	}

	if existingUser != nil {
		return nil, api_errors.ErrorUserAlreadyExists
	}

	newUser, newUserErr := s.userService.CreateUserByCredentials(ctx, dto.Email, dto.Name, dto.Password)

	if newUserErr != nil {
		s.logger.Error(fmt.Sprintf("failed to create user by credentials: %s", newUserErr))
		return nil, newUserErr
	}

	return newUser, nil
}

func (s *authorizationService) GetUserByToken(ctx context.Context, token string) (*common_models.User, error) {
	var authClaims common_types.AuthClaims

	err := s.jwt.Verify(token, &authClaims)

	if err != nil || authClaims.UserId == "" {
		return nil, api_errors.ErrorInvalidToken
	}

	user, userErr := s.userService.GetById(ctx, authClaims.UserId)

	if userErr != nil {
		s.logger.Error(fmt.Sprintf("failed to get user by id: %s", userErr))
		return nil, userErr
	}

	if user == nil {
		s.logger.Warn(fmt.Sprintf("user not found by id %d, but the token is valid %s", authClaims.UserId, token))
		return nil, api_errors.ErrorInvalidToken
	}

	return user, nil
}

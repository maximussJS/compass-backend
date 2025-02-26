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
	common_repositories "compass-backend/pkg/common/repositories"
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

	Logger         common_lib.ILogger
	Jwt            lib.IJwt
	Claims         lib.IClaims
	UserRepository common_repositories.IUserRepository
}

type authorizationService struct {
	logger         common_lib.ILogger
	jwt            lib.IJwt
	claims         lib.IClaims
	userRepository common_repositories.IUserRepository
}

func FxAuthorizationService() fx.Option {
	return fx_utils.AsProvider(newAuthorizationService, new(IAuthorizationService))
}

func newAuthorizationService(params authorizationServiceParams) *authorizationService {
	return &authorizationService{
		logger:         params.Logger,
		jwt:            params.Jwt,
		claims:         params.Claims,
		userRepository: params.UserRepository,
	}
}

func (s *authorizationService) SignInByPassword(ctx context.Context, dto authorization_dto.SignInByPasswordRequest) (string, error) {
	user, userErr := s.userRepository.GetByEmail(ctx, dto.Email)

	if userErr != nil {
		s.logger.Error(fmt.Sprintf("failed to get user by email: %s", userErr))
		return "", userErr
	}

	fmt.Printf("user: %t %s %t\n", user == nil, user.Password, crypto_utils.VerifyHash(dto.Password, user.Password))

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
	existingUser, existingErr := s.userRepository.GetByEmail(ctx, dto.Email)

	if existingErr != nil {
		s.logger.Error(fmt.Sprintf("failed to get user by email: %s", existingErr))
		return nil, existingErr
	}

	if existingUser != nil {
		return existingUser, api_errors.ErrorUserAlreadyExists
	}

	hashedPassword, hashErr := crypto_utils.Hash(dto.Password)

	if hashErr != nil {
		s.logger.Error(fmt.Sprintf("failed to hash password: %s", hashErr))
		return nil, hashErr
	}

	user := common_models.User{
		Email:    dto.Email,
		Name:     dto.Name,
		Password: hashedPassword,
	}

	id, createErr := s.userRepository.Create(ctx, user)

	if createErr != nil {
		s.logger.Error(fmt.Sprintf("failed to create user: %s", createErr))
		return nil, createErr
	}

	newUser, newUserErr := s.userRepository.GetById(ctx, id)

	if newUserErr != nil {
		s.logger.Error(fmt.Sprintf("failed to get new user by id: %s", newUserErr))
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

	user, userErr := s.userRepository.GetById(ctx, authClaims.UserId)

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

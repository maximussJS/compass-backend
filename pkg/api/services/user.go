package services

import (
	"compass-backend/pkg/api/api_errors"
	crypto_utils "compass-backend/pkg/api/utils/crypto"
	"compass-backend/pkg/api/utils/password"
	"compass-backend/pkg/common/constants"
	fx_utils "compass-backend/pkg/common/fx"
	common_lib "compass-backend/pkg/common/lib"
	"compass-backend/pkg/common/models"
	common_repositories "compass-backend/pkg/common/repositories"
	common_types "compass-backend/pkg/common/types"
	"context"
	"fmt"
	"go.uber.org/fx"
)

type IUserService interface {
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	CreateEmptyUserByEmail(ctx context.Context, email string) (*models.User, error)
	ChangeUserTeam(ctx context.Context, teamId string, user *models.User) error
}

type userServiceParams struct {
	fx.In

	Logger               common_lib.ILogger
	SenderService        ISenderService
	UserRepository       common_repositories.IUserRepository
	TeamMemberRepository common_repositories.ITeamMemberRepository
}

type userService struct {
	logger               common_lib.ILogger
	senderService        ISenderService
	userRepository       common_repositories.IUserRepository
	teamMemberRepository common_repositories.ITeamMemberRepository
}

func FxUserService() fx.Option {
	return fx_utils.AsProvider(newUserService, new(IUserService))
}

func newUserService(params userServiceParams) IUserService {
	return &userService{
		logger:               params.Logger,
		senderService:        params.SenderService,
		userRepository:       params.UserRepository,
		teamMemberRepository: params.TeamMemberRepository,
	}
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to get user by email: %s", err))
		return nil, err
	}

	return user, nil
}

func (s *userService) CreateEmptyUserByEmail(ctx context.Context, email string) (*models.User, error) {
	userPassword := password.GeneratePassword()

	hashedPassword, hashErr := crypto_utils.Hash(userPassword)

	if hashErr != nil {
		s.logger.Error(fmt.Sprintf("failed to hash password: %s", hashErr))
		return nil, hashErr
	}

	fmt.Printf("User password: %s\n", userPassword)
	fmt.Printf("Hashed password: %s\n", hashedPassword)

	id, createErr := s.userRepository.Create(ctx, models.User{
		Name:     "Unnamed User",
		Email:    email,
		Password: hashedPassword,
		Role:     constants.User,
	})

	if createErr != nil {
		s.logger.Error(fmt.Sprintf("failed to create user: %s", createErr))
		return nil, createErr
	}

	job := common_types.SendUserRegisteredEmailJob{
		Email:    email,
		Password: userPassword,
	}

	sendErr := s.senderService.SendUserRegistered(ctx, job)

	if sendErr != nil {
		s.logger.Error(fmt.Sprintf("failed to send user registered email: %s", sendErr))
		return nil, sendErr
	}

	user, getErr := s.userRepository.GetById(ctx, id)

	if getErr != nil {
		s.logger.Error(fmt.Sprintf("failed to get user: %s", getErr))
		return nil, getErr
	}

	return user, nil
}

func (s *userService) ChangeUserTeam(ctx context.Context, teamId string, user *models.User) error {
	if user.IsInTeam(teamId) {
		return api_errors.ErrorUserAlreadyInTeam
	}

	createErr := s.teamMemberRepository.Create(ctx, models.TeamMember{
		TeamId: teamId,
		UserId: user.Id,
	})

	if createErr != nil {
		s.logger.Error(fmt.Sprintf("failed to create team member: %s", createErr))
		return createErr
	}

	return nil
}

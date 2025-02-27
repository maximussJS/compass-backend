package services

import (
	"compass-backend/pkg/api/api_errors"
	user_dto "compass-backend/pkg/api/common/dto/user"
	crypto_utils "compass-backend/pkg/api/utils/crypto"
	"compass-backend/pkg/api/utils/password"
	"compass-backend/pkg/common/constants"
	fx_utils "compass-backend/pkg/common/fx"
	common_lib "compass-backend/pkg/common/lib"
	"compass-backend/pkg/common/models"
	common_repositories "compass-backend/pkg/common/repositories"
	common_services "compass-backend/pkg/common/services"
	common_types "compass-backend/pkg/common/types"
	"context"
	"fmt"
	"go.uber.org/fx"
)

type IUserService interface {
	ChangeName(ctx context.Context, user *models.User, request user_dto.ChangeNameRequest) error
	ChangePassword(ctx context.Context, user *models.User, request user_dto.ChangePasswordRequest) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetById(ctx context.Context, id string) (*models.User, error)
	CreateUserByCredentials(ctx context.Context, email, name, password string) (*models.User, error)
	CreateEmptyUserByEmail(ctx context.Context, email string) (*models.User, error)
	ChangeUserTeam(ctx context.Context, teamId string, user *models.User) error
}

type userServiceParams struct {
	fx.In

	Logger               common_lib.ILogger
	EmailSender          common_services.IEmailSenderService
	UserRepository       common_repositories.IUserRepository
	TeamMemberRepository common_repositories.ITeamMemberRepository
}

type userService struct {
	logger               common_lib.ILogger
	emailSender          common_services.IEmailSenderService
	userRepository       common_repositories.IUserRepository
	teamMemberRepository common_repositories.ITeamMemberRepository
}

func FxUserService() fx.Option {
	return fx_utils.AsProvider(newUserService, new(IUserService))
}

func newUserService(params userServiceParams) IUserService {
	return &userService{
		logger:               params.Logger,
		emailSender:          params.EmailSender,
		userRepository:       params.UserRepository,
		teamMemberRepository: params.TeamMemberRepository,
	}
}

func (s *userService) ChangeName(ctx context.Context, user *models.User, request user_dto.ChangeNameRequest) error {
	if user.Name == request.Name {
		return api_errors.ErrorNameIsTheSame
	}

	err := s.userRepository.UpdateById(ctx, user.Id, models.User{
		Name: request.Name,
	})

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to update user name: %s", err))
		return err
	}

	return nil
}

func (s *userService) ChangePassword(ctx context.Context, user *models.User, request user_dto.ChangePasswordRequest) error {
	if request.OldPassword == request.Password {
		return api_errors.ErrorUserPasswordIsTheSame
	}

	passwordIsValid := crypto_utils.VerifyHash(request.OldPassword, user.Password)

	if !passwordIsValid {
		return api_errors.ErrorInvalidCredentials
	}

	hashedPassword, hashErr := crypto_utils.Hash(request.Password)

	if hashErr != nil {
		s.logger.Error(fmt.Sprintf("failed to hash password: %s", hashErr))
		return hashErr
	}

	err := s.userRepository.UpdateById(ctx, user.Id, models.User{
		Password: hashedPassword,
	})

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to update user password: %s", err))
		return err
	}

	return nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to get user by email: %s", err))
		return nil, err
	}

	return user, nil
}

func (s *userService) GetById(ctx context.Context, id string) (*models.User, error) {
	user, err := s.userRepository.GetById(ctx, id)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to get user by id: %s", err))
		return nil, err
	}

	return user, nil
}

func (s *userService) CreateUserByCredentials(ctx context.Context, email, name, password string) (*models.User, error) {
	hashedPassword, hashErr := crypto_utils.Hash(password)

	if hashErr != nil {
		s.logger.Error(fmt.Sprintf("failed to hash password: %s", hashErr))
		return nil, hashErr
	}

	user := models.User{
		Email:      email,
		Name:       name,
		Password:   hashedPassword,
		IsVerified: false,
		Role:       constants.User,
	}

	id, createErr := s.userRepository.Create(ctx, user)

	if createErr != nil {
		s.logger.Error(fmt.Sprintf("failed to create user: %s", createErr))
		return nil, createErr
	}

	job := common_types.SendUserRegisteredEmailJobData{
		Name: name,
	}

	sendErr := s.emailSender.SendUserRegistered(ctx, job)

	if sendErr != nil {
		s.logger.Error(fmt.Sprintf("failed to send user registered email: %s", sendErr))
		return nil, sendErr
	}

	newUser, newUserErr := s.userRepository.GetById(ctx, id)

	if newUserErr != nil {
		s.logger.Error(fmt.Sprintf("failed to get new user by id: %s", newUserErr))
		return nil, newUserErr
	}

	return newUser, nil
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
		Name:       "Unnamed User",
		Email:      email,
		Password:   hashedPassword,
		IsVerified: true,
		Role:       constants.User,
	})

	if createErr != nil {
		s.logger.Error(fmt.Sprintf("failed to create user: %s", createErr))
		return nil, createErr
	}

	job := common_types.SendEmptyUserCreatedEmailJobData{
		Email:    email,
		Password: userPassword,
	}

	sendErr := s.emailSender.SendEmptyUserCreated(ctx, job)

	if sendErr != nil {
		s.logger.Error(fmt.Sprintf("failed to send empty user created email: %s", sendErr))
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

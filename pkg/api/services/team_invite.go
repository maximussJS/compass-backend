package services

import (
	"compass-backend/pkg/api/api_errors"
	"compass-backend/pkg/api/lib"
	fx_utils "compass-backend/pkg/common/fx"
	common_lib "compass-backend/pkg/common/lib"
	"compass-backend/pkg/common/models"
	common_repositories "compass-backend/pkg/common/repositories"
	common_types "compass-backend/pkg/common/types"
	"context"
	"fmt"
	"go.uber.org/fx"
	"time"
)

type ITeamInviteService interface {
	InviteByEmail(ctx context.Context, email, ownerId string) error
	AcceptInvite(ctx context.Context, token string) error
	CancelInvite(ctx context.Context, token string) error
}

type teamInviteServiceParams struct {
	fx.In

	Env                  lib.IEnv
	Logger               common_lib.ILogger
	TokenService         ITokenService
	EmailSender          IEmailSenderService
	UserService          IUserService
	TeamInviteRepository common_repositories.ITeamInviteRepository
	TeamRepository       common_repositories.ITeamRepository
}

type teamInviteService struct {
	appUrl               string
	logger               common_lib.ILogger
	tokenService         ITokenService
	emailSender          IEmailSenderService
	userService          IUserService
	teamInviteRepository common_repositories.ITeamInviteRepository
	teamRepository       common_repositories.ITeamRepository
}

func FxTeamInviteService() fx.Option {
	return fx_utils.AsProvider(newTeamInviteService, new(ITeamInviteService))
}

func newTeamInviteService(params teamInviteServiceParams) ITeamInviteService {
	return &teamInviteService{
		appUrl:               params.Env.GetAppUrl(),
		tokenService:         params.TokenService,
		logger:               params.Logger,
		userService:          params.UserService,
		emailSender:          params.EmailSender,
		teamRepository:       params.TeamRepository,
		teamInviteRepository: params.TeamInviteRepository,
	}
}

func (s *teamInviteService) InviteByEmail(ctx context.Context, email, ownerId string) error {
	team, teamErr := s.teamRepository.GetByOwnerId(ctx, ownerId)

	if teamErr != nil {
		s.logger.Error(fmt.Sprintf("failed to get team by owner id: %s", teamErr))
		return teamErr
	}

	if team == nil {
		return api_errors.ErrorTeamNotFound
	}

	existingTeamInvite, existingErr := s.teamInviteRepository.GetByEmailAndTeamId(ctx, email, team.Id)

	if existingErr != nil {
		s.logger.Error(fmt.Sprintf("failed to get team invite by email: %s", existingErr))
		return existingErr
	}

	if existingTeamInvite != nil {
		if existingTeamInvite.IsSent {
			return api_errors.ErrorTeamInviteAlreadySend
		}

		return nil
	}

	token, expiresAt, tokenErr := s.tokenService.GenerateTeamInviteToken(team.Id, email)

	if tokenErr != nil {
		s.logger.Error(fmt.Sprintf("failed to generate team invite token: %s", tokenErr))
		return tokenErr
	}

	inviteId, inviteErr := s.teamInviteRepository.Create(ctx, models.TeamInvite{
		Email:     email,
		TeamId:    team.Id,
		Token:     token,
		ExpiresAt: time.Unix(expiresAt, 0),
	})

	if inviteErr != nil {
		s.logger.Error(fmt.Sprintf("failed to create team invite: %s", inviteErr))
		return inviteErr
	}

	job := common_types.SendTeamInviteEmailJobData{
		Id:         inviteId,
		AcceptLink: fmt.Sprintf("%s/api/team-invites/accept/%s", s.appUrl, token),
		CancelLink: fmt.Sprintf("%s/api/team-invites/cancel/%s", s.appUrl, token),
	}

	sendErr := s.emailSender.SendTeamInvite(ctx, job)

	if sendErr != nil {
		s.logger.Error(fmt.Sprintf("failed to send team invite: %s", sendErr))
		return sendErr
	}

	return nil
}

func (s *teamInviteService) AcceptInvite(ctx context.Context, token string) error {
	invite, inviteErr := s.verifyTeamInviteByToken(ctx, token)

	if inviteErr != nil {
		return inviteErr
	}

	existingUser, existingErr := s.userService.GetByEmail(ctx, invite.Email)

	if existingErr != nil {
		return existingErr
	}

	if existingUser != nil {
		return s.changeUserTeam(ctx, invite, existingUser)
	}

	newUser, newUserErr := s.userService.CreateEmptyUserByEmail(ctx, invite.Email)

	if newUserErr != nil {
		s.logger.Error(fmt.Sprintf("failed to create user by email: %s", newUserErr))
		return newUserErr
	}

	return s.changeUserTeam(ctx, invite, newUser)
}

func (s *teamInviteService) changeUserTeam(ctx context.Context, invite *models.TeamInvite, user *models.User) error {
	changeErr := s.userService.ChangeUserTeam(ctx, invite.TeamId, user)

	if changeErr != nil {
		s.logger.Error(fmt.Sprintf("failed to change user team: %s", changeErr))
		return changeErr
	}

	return s.teamInviteRepository.MarkAsAccepted(ctx, invite.Id)
}

func (s *teamInviteService) verifyTeamInviteByToken(ctx context.Context, token string) (*models.TeamInvite, error) {
	inviteClaims, inviteClaimsErr := s.tokenService.VerifyTeamInviteToken(token)

	if inviteClaimsErr != nil {
		return nil, inviteClaimsErr
	}

	if inviteClaims.IsExpired() {
		return nil, api_errors.ErrorTeamInviteExpired
	}

	invite, inviteErr := s.teamInviteRepository.GetByToken(ctx, token)

	if inviteErr != nil {
		s.logger.Error(fmt.Sprintf("failed to get team invite by token: %s", inviteErr))
		return nil, inviteErr
	}

	if invite == nil {
		return nil, api_errors.ErrorTeamInviteNotFound
	}

	if invite.IsExpired() {
		return nil, api_errors.ErrorTeamInviteExpired
	}

	if invite.IsAccepted() {
		return nil, api_errors.ErrorTeamInviteAccepted
	}

	if invite.IsCancelled() {
		return nil, api_errors.ErrorTeamInviteCancelled
	}

	if invite.Email != inviteClaims.Email {
		s.logger.Warn(fmt.Sprintf("email mismatch in invite: %s != %s", invite.Email, inviteClaims.Email))
		return nil, api_errors.ErrorInvalidToken
	}

	return invite, nil
}

func (s *teamInviteService) CancelInvite(ctx context.Context, token string) error {
	invite, inviteErr := s.verifyTeamInviteByToken(ctx, token)

	if inviteErr != nil {
		return inviteErr
	}

	return s.teamInviteRepository.MarkAsCancelled(ctx, invite.Id)
}

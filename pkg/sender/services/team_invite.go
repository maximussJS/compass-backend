package services

import (
	fx_utils "compass-backend/pkg/common/fx"
	"compass-backend/pkg/common/models"
	common_repositories "compass-backend/pkg/common/repositories"
	"compass-backend/pkg/sender/lib"
	"compass-backend/pkg/sender/template_data"
	"context"
	"fmt"
	"go.uber.org/fx"
)

type ITeamInviteService interface {
	SendTeamInvite(ctx context.Context, teamInviteId, acceptLink, cancelLink string) error
}

type teamInviteServiceParams struct {
	fx.In

	HtmlTemplate         lib.IHtmlTemplate
	MailService          IMailService
	TeamInviteRepository common_repositories.ITeamInviteRepository
}

type teamInviteService struct {
	htmlTemplate         lib.IHtmlTemplate
	mailService          IMailService
	teamInviteRepository common_repositories.ITeamInviteRepository
}

func FxTeamInviteService() fx.Option {
	return fx_utils.AsProvider(newTeamInviteService, new(ITeamInviteService))
}

func newTeamInviteService(params teamInviteServiceParams) ITeamInviteService {
	return &teamInviteService{
		htmlTemplate:         params.HtmlTemplate,
		mailService:          params.MailService,
		teamInviteRepository: params.TeamInviteRepository,
	}
}

func (s *teamInviteService) SendTeamInvite(ctx context.Context, teamInviteId, acceptLink, cancelLink string) error {
	teamInvite, err := s.teamInviteRepository.GetById(ctx, teamInviteId)
	if err != nil {
		return fmt.Errorf("failed to get team invite by id: %v", err)
	}

	if teamInvite == nil {
		return fmt.Errorf("team invite %s not found", teamInviteId)
	}

	data := template_data.TeamInviteTemplateData{
		AcceptLink: acceptLink,
		CancelLink: cancelLink,
		OwnerEmail: teamInvite.TeamOwnerEmail(),
		TeamName:   teamInvite.TeamName(),
		OwnerName:  teamInvite.TeamOwnerName(),
		ExpiresAt:  teamInvite.ExpiresAt.Format("02 Jan 2006 15:04:05"),
	}

	template, err := s.htmlTemplate.NewTeamInviteTemplate(data)

	if err != nil {
		return fmt.Errorf("failed to create team invite template: %v", err)
	}

	sendErr := s.mailService.Send(teamInvite.Email, "Compass App Team Invitation", template)

	if sendErr != nil {
		return fmt.Errorf("failed to send team invite email: %v", err)
	}

	updateErr := s.teamInviteRepository.UpdateById(ctx, teamInviteId, models.TeamInvite{
		IsSent: true,
	})

	if updateErr != nil {
		return fmt.Errorf("failed to update team invite: %v", err)
	}

	return nil
}

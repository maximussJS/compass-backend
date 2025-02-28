package services

import (
	fx_utils "compass-backend/pkg/common/fx"
	"compass-backend/pkg/sender/lib"
	"compass-backend/pkg/sender/template_data"
	"context"
	"fmt"
	"go.uber.org/fx"
)

type IUserService interface {
	SendEmptyUserCreated(ctx context.Context, email string, password string) error
	SendConfirmEmail(ctx context.Context, email, name, confirmationLink string) error
}

type userServiceParams struct {
	fx.In

	HtmlTemplate lib.IHtmlTemplate
	MailService  IMailService
}

type userService struct {
	htmlTemplate lib.IHtmlTemplate
	mailService  IMailService
}

func FxUserService() fx.Option {
	return fx_utils.AsProvider(newUserService, new(IUserService))
}

func newUserService(params userServiceParams) IUserService {
	return &userService{
		htmlTemplate: params.HtmlTemplate,
		mailService:  params.MailService,
	}
}

func (s *userService) SendEmptyUserCreated(_ context.Context, email string, password string) error {
	data := template_data.EmptyUserCreatedTemplateData{
		Email:    email,
		Password: password,
	}

	template, err := s.htmlTemplate.NewUserCreatedTemplate(data)

	if err != nil {
		return fmt.Errorf("failed to create empty user created template: %v", err)
	}

	sendErr := s.mailService.Send(email, "Compass App Registration", template)

	if sendErr != nil {
		return fmt.Errorf("failed to send empty user created email: %v", err)
	}

	return nil
}

func (s *userService) SendConfirmEmail(_ context.Context, email, name, confirmationLink string) error {
	data := template_data.ConfirmEmailTemplateData{
		Name:             name,
		ConfirmationLink: confirmationLink,
	}

	template, err := s.htmlTemplate.NewConfirmEmailTemplate(data)

	if err != nil {
		return fmt.Errorf("failed to create confirm email template: %v", err)
	}

	sendErr := s.mailService.Send(email, "Compass App Registration", template)

	if sendErr != nil {
		return fmt.Errorf("failed to send confirm email: %v", err)
	}

	return nil
}

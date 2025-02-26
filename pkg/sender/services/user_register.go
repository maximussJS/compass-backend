package services

import (
	fx_utils "compass-backend/pkg/common/fx"
	"compass-backend/pkg/sender/lib"
	"compass-backend/pkg/sender/template_data"
	"context"
	"fmt"
	"go.uber.org/fx"
)

type IUserRegisterService interface {
	SendRegisterUser(ctx context.Context, email string, password string) error
}

type userRegisterServiceParams struct {
	fx.In

	HtmlTemplate lib.IHtmlTemplate
	MailService  IMailService
}

type userRegisterService struct {
	htmlTemplate lib.IHtmlTemplate
	mailService  IMailService
}

func FxUserRegisterService() fx.Option {
	return fx_utils.AsProvider(newUserRegisterService, new(IUserRegisterService))
}

func newUserRegisterService(params userRegisterServiceParams) IUserRegisterService {
	return &userRegisterService{
		htmlTemplate: params.HtmlTemplate,
		mailService:  params.MailService,
	}
}

func (s *userRegisterService) SendRegisterUser(ctx context.Context, email string, password string) error {
	data := template_data.UserRegisterTemplateData{
		Email:    email,
		Password: password,
	}

	template, err := s.htmlTemplate.NewUserRegisterTemplate(data)

	if err != nil {
		return fmt.Errorf("failed to create user register template: %v", err)
	}

	sendErr := s.mailService.Send(email, "Compass App Registration", template)

	if sendErr != nil {
		return fmt.Errorf("failed to send user register email: %v", err)
	}

	return nil
}

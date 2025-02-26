package lib

import (
	"bytes"
	fx_utils "compass-backend/pkg/common/fx"
	"compass-backend/pkg/sender/template_data"
	"fmt"
	"go.uber.org/fx"
	"html/template"
)

type IHtmlTemplate interface {
	NewTeamInviteTemplate(data template_data.TeamInviteTemplateData) (string, error)
	NewUserRegisterTemplate(data template_data.UserRegisterTemplateData) (string, error)
}

type htmlTemplateParams struct {
	fx.In

	Env IEnv
}

type htmlTemplate struct {
	templateDirectory string
}

func FxHtmlTemplate() fx.Option {
	return fx_utils.AsProvider(newHtmlTemplate, new(IHtmlTemplate))
}

func newHtmlTemplate(params htmlTemplateParams) IHtmlTemplate {
	return &htmlTemplate{
		templateDirectory: params.Env.GetHtmlTemplateDirectory(),
	}
}

func (h *htmlTemplate) NewTeamInviteTemplate(data template_data.TeamInviteTemplateData) (string, error) {
	filepath := fmt.Sprintf("%s/invite.html", h.templateDirectory)

	t, err := template.New("invite.html").ParseFiles(filepath)

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func (h *htmlTemplate) NewUserRegisterTemplate(data template_data.UserRegisterTemplateData) (string, error) {
	filepath := fmt.Sprintf("%s/user_registered.html", h.templateDirectory)

	t, err := template.New("user_registered.html").ParseFiles(filepath)

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

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
	NewUserCreatedTemplate(data template_data.EmptyUserCreatedTemplateData) (string, error)
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

func (h *htmlTemplate) NewUserCreatedTemplate(data template_data.EmptyUserCreatedTemplateData) (string, error) {
	filepath := fmt.Sprintf("%s/new_user_created.html", h.templateDirectory)

	t, err := template.New("new_user_created.html").ParseFiles(filepath)

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

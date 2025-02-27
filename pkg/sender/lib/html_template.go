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
	return h.createTemplate("team_invite", data)
}

func (h *htmlTemplate) NewUserCreatedTemplate(data template_data.EmptyUserCreatedTemplateData) (string, error) {
	return h.createTemplate("new_user_created", data)
}

func (h *htmlTemplate) createTemplate(templateName string, data interface{}) (string, error) {
	filepath := fmt.Sprintf("%s/%s.html", h.templateDirectory, templateName)

	t, err := template.New(fmt.Sprintf("%s.html", templateName)).ParseFiles(filepath)

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

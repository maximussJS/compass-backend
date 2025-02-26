package lib

import (
	fx_utils "compass-backend/pkg/common/fx"
	common_lib "compass-backend/pkg/common/lib"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type IEnv interface {
	common_lib.IEnv
	GetMailUsername() string
	GetMailPassword() string
	GetMailHost() string
	GetMailPort() int
	GetHtmlTemplateDirectory() string
	GetTeamInviteRedisChannel() string
	GetUserRegisteredChannel() string
}

type envParams struct {
	fx.In

	common_lib.IEnv
}

type env struct {
	common_lib.IEnv
	MailUsername           string `mapstructure:"MAIL_USERNAME" validate:"required"`
	MailPassword           string `mapstructure:"MAIL_PASSWORD" validate:"required"`
	MailHost               string `mapstructure:"MAIL_HOST" validate:"required"`
	MailPort               int    `mapstructure:"MAIL_PORT" validate:"required"`
	HtmlTemplateDirectory  string `mapstructure:"HTML_TEMPLATE_DIRECTORY"`
	TeamInviteRedisChannel string `mapstructure:"TEAM_INVITE_REDIS_CHANNEL"`
	UserRegisteredChannel  string `mapstructure:"USER_REGISTERED_REDIS_CHANNEL"`
}

func FxEnv() fx.Option {
	return fx_utils.AsProvider(newEnv, new(IEnv))
}

func newEnv(params envParams) *env {
	_ = godotenv.Load()

	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		panic("cannot read configuration")
	}

	viper.SetDefault("TIMEZONE", "UTC")

	env := &env{
		IEnv:                   params.IEnv,
		HtmlTemplateDirectory:  "pkg/sender/templates",
		TeamInviteRedisChannel: "team_invite",
		UserRegisteredChannel:  "user_registered",
	}

	err = viper.Unmarshal(env)
	if err != nil {
		panic(fmt.Sprintf("environment cant be loaded: %v", err))
	}

	validate := validator.New()
	if err := validate.Struct(env); err != nil {
		panic(fmt.Sprintf("environment validation failed: %v", err))
	}

	return env
}

func (e *env) GetMailUsername() string {
	return e.MailUsername
}

func (e *env) GetMailPassword() string {
	return e.MailPassword
}

func (e *env) GetMailHost() string {
	return e.MailHost
}

func (e *env) GetMailPort() int {
	return e.MailPort
}

func (e *env) GetHtmlTemplateDirectory() string {
	return e.HtmlTemplateDirectory
}

func (e *env) GetTeamInviteRedisChannel() string {
	return e.TeamInviteRedisChannel
}

func (e *env) GetUserRegisteredChannel() string {
	return e.UserRegisteredChannel
}

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
	GetEmailRedisChannel() string
}

type envParams struct {
	fx.In

	common_lib.IEnv
}

type env struct {
	common_lib.IEnv
	MailUsername          string `mapstructure:"MAIL_USERNAME" validate:"required"`
	MailPassword          string `mapstructure:"MAIL_PASSWORD" validate:"required"`
	MailHost              string `mapstructure:"MAIL_HOST" validate:"required"`
	MailPort              int    `mapstructure:"MAIL_PORT" validate:"required"`
	EmailRedisChannel     string `mapstructure:"EMAIL_REDIS_CHANNEL"`
	HtmlTemplateDirectory string `mapstructure:"HTML_TEMPLATE_DIRECTORY"`
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
		IEnv:                  params.IEnv,
		HtmlTemplateDirectory: "pkg/sender/templates",
		EmailRedisChannel:     "email-channel",
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

func (e *env) GetEmailRedisChannel() string {
	return e.EmailRedisChannel
}

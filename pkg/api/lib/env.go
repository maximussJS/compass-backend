package lib

import (
	fx_utils "compass-backend/pkg/common/fx"
	common_lib "compass-backend/pkg/common/lib"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"time"
)

type IEnv interface {
	common_lib.IEnv
	GetAppUrl() string
	GetDefaultLimit() int
	GetMaxLimit() int
	GetJwtSecret() []byte
	GetAuthExpirationDuration() time.Duration
	GetInviteExpirationDuration() time.Duration
	GetEmailRedisChannel() string
}

type envParams struct {
	fx.In

	common_lib.IEnv
}

type env struct {
	common_lib.IEnv
	JwtSecret string `mapstructure:"JWT_SECRET" validate:"required"`
	AppUrl    string `mapstructure:"APP_URL" validate:"required"`

	AuthExpirationDurationInMinutes   int `mapstructure:"AUTH_EXPIRATION_DURATION_IN_MINUTES"`
	InviteExpirationDurationInMinutes int `mapstructure:"INVITE_EXPIRATION_DURATION_IN_MINUTES"`

	EmailRedisChannel string `mapstructure:"EMAIL_REDIS_CHANNEL"`

	DefaultLimit int `mapstructure:"DEFAULT_LIMIT"`
	MaxLimit     int `mapstructure:"MAX_LIMIT"`
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
		IEnv:                              params.IEnv,
		DefaultLimit:                      10,
		MaxLimit:                          100,
		AuthExpirationDurationInMinutes:   1440, // 24 hours
		InviteExpirationDurationInMinutes: 1,
		EmailRedisChannel:                 "email-channel",
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

func (e *env) GetDefaultLimit() int {
	return e.DefaultLimit
}

func (e *env) GetMaxLimit() int {
	return e.MaxLimit
}

func (e *env) GetJwtSecret() []byte {
	return []byte(e.JwtSecret)
}

func (e *env) GetAuthExpirationDuration() time.Duration {
	return time.Duration(e.AuthExpirationDurationInMinutes) * time.Minute
}

func (e *env) GetInviteExpirationDuration() time.Duration {
	return time.Duration(e.InviteExpirationDurationInMinutes) * time.Minute
}

func (e *env) GetAppUrl() string {
	return e.AppUrl
}

func (e *env) GetEmailRedisChannel() string {
	return e.EmailRedisChannel
}

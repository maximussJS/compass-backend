package lib

import (
	"compass-backend/pkg/common/constants"
	fx_utils "compass-backend/pkg/common/fx"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"
	"time"
)

type IEnv interface {
	GetAppName() string
	GetPostgresUrl() string
	GetRedisUrl() string
	GetLoggerLevel() zapcore.Level
	GetEnvironment() constants.AppEnv
	GetTimeZone() string
	GetPort() int
	GetRequestTimeoutDuration() time.Duration
	GetMaxMultipartMemory() int64
}

type env struct {
	AppName                 string `mapstructure:"APP_NAME" validate:"required"`
	PostgresUrl             string `mapstructure:"POSTGRES_URL" validate:"required"`
	RedisUrl                string `mapstructure:"REDIS_URL" validate:"required"`
	LoggerLevel             string `mapstructure:"LOGGER_LEVEL" validate:"required"`
	Environment             string `mapstructure:"ENVIRONMENT" validate:"required"`
	Port                    int    `mapstructure:"PORT" validate:"required"`
	TimeZone                string `mapstructure:"TIMEZONE"`
	RequestTimeoutInSeconds int    `mapstructure:"REQUEST_TIMEOUT_IN_SECONDS"`
	MaxMultipartMemory      int64  `mapstructure:"MAX_MULTIPART_MEMORY"`
}

func FxEnv() fx.Option {
	return fx_utils.AsProvider(newEnv, new(IEnv))
}

func newEnv() *env {
	_ = godotenv.Load()

	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		panic("cannot read configuration")
	}

	viper.SetDefault("TIMEZONE", "UTC")

	env := &env{
		RequestTimeoutInSeconds: 30,
		MaxMultipartMemory:      10 << 20, // 10 MB
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

func (e *env) GetLoggerLevel() zapcore.Level {
	switch e.LoggerLevel {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.PanicLevel
	}
}

func (e *env) GetEnvironment() constants.AppEnv {
	switch e.Environment {
	case "development":
		return constants.DevelopmentEnv
	case "production":
		return constants.ProductionEnv
	default:
		panic(fmt.Sprintf("unknown environment: %s", e.Environment))
	}
}

func (e *env) GetTimeZone() string {
	return e.TimeZone
}

func (e *env) GetPort() int {
	return e.Port
}

func (e *env) GetRedisUrl() string {
	return e.RedisUrl
}

func (e *env) GetPostgresUrl() string {
	return e.PostgresUrl
}

func (e *env) GetRequestTimeoutDuration() time.Duration {
	return time.Duration(e.RequestTimeoutInSeconds) * time.Second
}

func (e *env) GetAppName() string {
	return e.AppName
}

func (e *env) GetMaxMultipartMemory() int64 {
	return e.MaxMultipartMemory
}

package lib

import (
	"compass-backend/pkg/common/constants"
	fx_utils "compass-backend/pkg/common/fx"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ILogger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
}

type logger struct {
	*zap.Logger
}

type loggerParams struct {
	fx.In

	Env IEnv
}

func FxLogger() fx.Option {
	return fx_utils.AsProvider(newLogger, new(ILogger))
}

func newLogger(params loggerParams) *logger {
	config := zap.NewDevelopmentConfig()

	if params.Env.GetEnvironment() == constants.DevelopmentEnv {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.Level.SetLevel(params.Env.GetLoggerLevel())

	zapLogger, _ := config.Build()

	return &logger{
		zapLogger,
	}
}

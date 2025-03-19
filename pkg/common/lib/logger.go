package lib

import (
	"compass-backend/pkg/common/constants"
	fx_utils "compass-backend/pkg/common/fx"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ILogger interface {
	Info(msg string, fields ...zap.Field)
	Infof(msg string, args ...interface{})
	Error(msg string, fields ...zap.Field)
	Errorf(msg string, args ...interface{})
	Fatal(msg string, fields ...zap.Field)
	Fatalf(msg string, args ...interface{})
	Warn(msg string, fields ...zap.Field)
	Warnf(msg string, args ...interface{})
	Debug(msg string, fields ...zap.Field)
	Debugf(msg string, args ...interface{})
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

func (l *logger) Errorf(msg string, args ...interface{}) {
	l.Error(fmt.Sprintf(msg, args...))
}

func (l *logger) Fatalf(msg string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(msg, args...))
}

func (l *logger) Warnf(msg string, args ...interface{}) {
	l.Warn(fmt.Sprintf(msg, args...))
}

func (l *logger) Debugf(msg string, args ...interface{}) {
	l.Debug(fmt.Sprintf(msg, args...))
}

func (l *logger) Infof(msg string, args ...interface{}) {
	l.Info(fmt.Sprintf(msg, args...))
}

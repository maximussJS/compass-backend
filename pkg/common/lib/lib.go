package lib

import (
	"compass-backend/pkg/common/constants"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"io"
	"os"
)

var Module = fx.Options(
	FxEnv(),
	FxLogger(),
	fx.WithLogger(func(env IEnv) fxevent.Logger {
		if env.GetEnvironment() == constants.ProductionEnv {
			return fxevent.NopLogger
		}
		return &fxevent.ConsoleLogger{
			W: io.Writer(os.Stdout),
		}
	}),
)

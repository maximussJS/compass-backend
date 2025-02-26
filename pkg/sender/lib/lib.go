package lib

import (
	common_lib "compass-backend/pkg/common/lib"
	"go.uber.org/fx"
)

var Module = fx.Options(
	common_lib.Module,
	FxEnv(),
	FxHtmlTemplate(),
)

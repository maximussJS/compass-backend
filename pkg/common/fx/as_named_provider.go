package fx

import "go.uber.org/fx"

func AsNamedProvider(provider any, providerInterface interface{}, tags string) fx.Option {
	return fx.Provide(fx.Annotate(provider, fx.As(providerInterface), fx.ResultTags(tags)))
}

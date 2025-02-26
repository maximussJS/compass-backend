package infrastructure

import (
	fx_utils "compass-backend/pkg/common/fx"
	"compass-backend/pkg/common/lib"
	"context"
	"errors"
	"fmt"
	"go.uber.org/fx"
	"net/http"
)

type IHttpServer interface {
}

type httpServer struct {
	*http.Server
	logger lib.ILogger
}

type httpServerParams struct {
	fx.In

	Router IRouter
	Env    lib.IEnv
	Logger lib.ILogger
}

func FxHttpServer() fx.Option {
	return fx_utils.AsProvider(newHttpServer, new(IHttpServer))
}

func newHttpServer(lc fx.Lifecycle, params httpServerParams) *httpServer {
	router := params.Router.GetRouter()

	server := &httpServer{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", params.Env.GetPort()),
			Handler: router,
		},
		logger: params.Logger,
	}

	lc.Append(fx.Hook{
		OnStart: server.Start,
		OnStop:  server.Stop,
	})

	return server
}

func (s *httpServer) Start(ctx context.Context) error {
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error(fmt.Sprintf("http server ListenAndServe error: %v", err))
		}
	}()

	s.logger.Info(fmt.Sprintf("http server started on port %s", s.Addr))
	return nil
}

func (s *httpServer) Stop(ctx context.Context) error {
	err := s.Shutdown(ctx)

	if err != nil {
		s.logger.Error(fmt.Sprintf("http server shutdown error: %v", err))
		return err
	}

	s.logger.Info("http server stopped")

	return nil
}

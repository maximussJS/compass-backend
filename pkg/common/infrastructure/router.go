package infrastructure

import (
	"compass-backend/pkg/common/api/middlewares"
	common_constants "compass-backend/pkg/common/constants"
	fx_utils "compass-backend/pkg/common/fx"
	"compass-backend/pkg/common/lib"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"time"
)

type IRouter interface {
	GetRouter() *gin.Engine
	Use(middlewares ...gin.HandlerFunc)
}

type router struct {
	httpRouter *gin.Engine
}

type routerParams struct {
	fx.In

	Env               lib.IEnv
	Logger            lib.ILogger
	TimeoutMiddleware middlewares.ITimeoutMiddleware
}

func FxRouter() fx.Option {
	return fx_utils.AsProvider(newRouter, new(IRouter))
}

func newRouter(params routerParams) *router {
	if params.Env.GetEnvironment() == common_constants.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	httpRouter := gin.New()

	err := httpRouter.SetTrustedProxies(nil)
	if err != nil {
		params.Logger.Fatal(err.Error())
	}

	httpRouter.MaxMultipartMemory = params.Env.GetMaxMultipartMemory()

	httpRouter.Use(ginzap.Ginzap(params.Logger, time.RFC3339, true))
	httpRouter.Use(ginzap.RecoveryWithZap(params.Logger, true))

	httpRouter.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	httpRouter.Use(params.TimeoutMiddleware.Handle())

	return &router{
		httpRouter: httpRouter,
	}
}

func (r *router) GetRouter() *gin.Engine {
	return r.httpRouter
}

func (r *router) Use(middlewares ...gin.HandlerFunc) {
	r.httpRouter.Use(middlewares...)
}

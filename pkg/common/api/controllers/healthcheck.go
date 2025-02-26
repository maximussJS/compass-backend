package controllers

import (
	fx_utils "compass-backend/pkg/common/fx"
	"compass-backend/pkg/common/lib"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
)

type IHealthCheckController interface {
	HealthCheck(c *gin.Context)
}

type healthCheckControllerParams struct {
	fx.In

	Env lib.IEnv
}

type healthCheckController struct {
	appName string
}

func FxHealthCheckController() fx.Option {
	return fx_utils.AsProvider(newHealthCheckController, new(IHealthCheckController))
}

func newHealthCheckController(params healthCheckControllerParams) IHealthCheckController {
	return &healthCheckController{
		appName: params.Env.GetAppName(),
	}
}

func (h *healthCheckController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": fmt.Sprintf("%s is up and running", h.appName)})
}

package routes

import (
	"compass-backend/pkg/api/api/controllers"
	"compass-backend/pkg/api/api/middlewares"
	common_routes "compass-backend/pkg/common/api/routes"
	common_infrastracture "compass-backend/pkg/common/infrastructure"
	common_lib "compass-backend/pkg/common/lib"
	"fmt"
	"go.uber.org/fx"
)

type categoryRoute struct {
	logger                  common_lib.ILogger
	authorizationMiddleware middlewares.IAuthorizationMiddleware
	router                  common_infrastracture.IRouter
	categoryController      controllers.ICategoryController
}

type categoryRouteParams struct {
	fx.In

	Logger                  common_lib.ILogger
	AuthorizationMiddleware middlewares.IAuthorizationMiddleware
	CategoryController      controllers.ICategoryController
	Router                  common_infrastracture.IRouter
}

func FxCategoryRoute() fx.Option {
	return common_routes.AsRoute(newCategoryRoute)
}

func newCategoryRoute(params categoryRouteParams) common_routes.IRoute {
	return &categoryRoute{
		logger:                  params.Logger,
		router:                  params.Router,
		authorizationMiddleware: params.AuthorizationMiddleware,
		categoryController:      params.CategoryController,
	}
}

func (h *categoryRoute) Setup() {
	group := h.router.GetRouter().Group("api/categories")

	group.Use(h.authorizationMiddleware.Handle())

	h.logger.Info(fmt.Sprintf("Mapped Category Route %s", group.BasePath()))

	group.POST("", h.categoryController.Create)
	h.logger.Info(fmt.Sprintf("Mapped POST %s", group.BasePath()))

	group.GET("", h.categoryController.List)
	h.logger.Info(fmt.Sprintf("Mapped GET %s", group.BasePath()))

	group.GET("/:id", h.categoryController.GetById)
	h.logger.Info(fmt.Sprintf("Mapped GET %s/:id", group.BasePath()))

	group.PATCH("/:id", h.categoryController.UpdateById)
	h.logger.Info(fmt.Sprintf("Mapped PATCH %s/:id", group.BasePath()))

	group.DELETE("/:id", h.categoryController.DeleteById)
	h.logger.Info(fmt.Sprintf("Mapped DELETE %s/:id", group.BasePath()))
}

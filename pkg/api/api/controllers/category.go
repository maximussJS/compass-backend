package controllers

import (
	"compass-backend/pkg/api/api/responses"
	"compass-backend/pkg/api/api_errors"
	category_dto "compass-backend/pkg/api/common/dto/category"
	"compass-backend/pkg/api/lib"
	"compass-backend/pkg/api/services"
	gin_utils "compass-backend/pkg/api/utils/gin"
	fx_utils "compass-backend/pkg/common/fx"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ICategoryController interface {
	Create(c *gin.Context)
	GetById(c *gin.Context)
	List(c *gin.Context)
	UpdateById(c *gin.Context)
	DeleteById(c *gin.Context)
}

type categoryControllerParams struct {
	fx.In

	Env             lib.IEnv
	CategoryService services.ICategoryService
}

type categoryController struct {
	defaultLimit int
	maxLimit     int
	service      services.ICategoryService
}

func FxCategoryController() fx.Option {
	return fx_utils.AsProvider(newCategoryController, new(ICategoryController))
}

func newCategoryController(params categoryControllerParams) ICategoryController {
	return &categoryController{
		defaultLimit: params.Env.GetDefaultLimit(),
		maxLimit:     params.Env.GetMaxLimit(),
		service:      params.CategoryService,
	}
}

func (h *categoryController) Create(c *gin.Context) {
	var req category_dto.CreateCategoryRequest

	if ok := gin_utils.BindData(c, &req); !ok {
		return
	}

	category, err := h.service.Create(c, req)

	if err != nil {
		if errors.Is(err, api_errors.ErrorCategoryAlreadyExists) {
			responses.BadRequest(c, fmt.Sprintf("Category with name %s already exists", req.Name))
			return
		}

		responses.InternalServerError(c)
		return
	}

	responses.SuccessJsonWithMessage(c, category, fmt.Sprintf("Category %s created", category.Name))
}

func (h *categoryController) GetById(c *gin.Context) {
	id, ok := gin_utils.GetUintIdParam(c)
	if !ok {
		return
	}

	category, err := h.service.GetById(c, id)

	if err != nil {
		if errors.Is(err, api_errors.ErrorCategoryNotFound) {
			responses.NotFound(c, fmt.Sprintf("Category with id %d not found", id))
			return
		}

		responses.InternalServerError(c)
		return
	}

	responses.SuccessJson(c, category)
}

func (h *categoryController) List(c *gin.Context) {
	limit := gin_utils.GetLimit(c, h.defaultLimit, h.maxLimit)
	offset := gin_utils.GetOffset(c)

	if limit == -1 || offset == -1 {
		return
	}

	categories, err := h.service.List(c, limit, offset)

	if err != nil {
		responses.InternalServerError(c)
		return
	}

	responses.SuccessJson(c, categories)
}

func (h *categoryController) UpdateById(c *gin.Context) {
	id, ok := gin_utils.GetUintIdParam(c)
	if !ok {
		return
	}

	var req category_dto.UpdateCategoryRequest

	if ok := gin_utils.BindData(c, &req); !ok {
		return
	}

	category, err := h.service.UpdateById(c, id, req)

	if err != nil {
		if errors.Is(err, api_errors.ErrorCategoryNotFound) {
			responses.NotFound(c, fmt.Sprintf("Category with id %d not found", id))
			return
		}

		if errors.Is(err, api_errors.ErrorCategoryAlreadyExists) {
			responses.BadRequest(c, fmt.Sprintf("Category with name %s already exists", req.Name))
			return
		}

		responses.InternalServerError(c)
		return
	}

	responses.SuccessJson(c, category)
}

func (h *categoryController) DeleteById(c *gin.Context) {
	id, ok := gin_utils.GetUintIdParam(c)
	if !ok {
		return
	}

	category, err := h.service.DeleteById(c, id)

	if err != nil {
		if errors.Is(err, api_errors.ErrorCategoryNotFound) {
			responses.NotFound(c, fmt.Sprintf("Category with id %d not found", id))
			return
		}

		responses.InternalServerError(c)
		return
	}

	responses.SuccessJson(c, category)
}

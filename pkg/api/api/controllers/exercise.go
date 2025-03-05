package controllers

import (
	"compass-backend/pkg/api/api/responses"
	"compass-backend/pkg/api/api_errors"
	dto_exercise "compass-backend/pkg/api/common/dto/exercise"
	"compass-backend/pkg/api/services"
	gin_utils "compass-backend/pkg/api/utils/gin"
	fx_utils "compass-backend/pkg/common/fx"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type IExerciseController interface {
	Create(c *gin.Context)
}

type exerciseControllerParams struct {
	fx.In

	ExerciseService services.IExerciseService
}

type exerciseController struct {
	exerciseService services.IExerciseService
}

func FxExerciseController() fx.Option {
	return fx_utils.AsProvider(newExerciseController, new(IExerciseController))
}

func newExerciseController(params exerciseControllerParams) IExerciseController {
	return &exerciseController{
		exerciseService: params.ExerciseService,
	}
}

func (h *exerciseController) Create(c *gin.Context) {
	user, userErr := gin_utils.GetUser(c)

	if userErr != nil {
		responses.UnauthorizedWithAbort(c)
		return
	}

	var data dto_exercise.CreateExerciseRequest

	if ok := gin_utils.BindData(c, &data); !ok {
		return
	}

	exercise, err := h.exerciseService.Create(c, user.Id, data)

	if err != nil {
		if errors.Is(err, api_errors.ErrorExerciseAlreadyExists) {
			responses.BadRequestWithAbort(c, fmt.Sprintf("Exercise with name %s already exists", data.Name))
			return
		}

		responses.InternalServerError(c)
		return
	}

	responses.SuccessJsonWithMessage(c, exercise, fmt.Sprintf("Exercise %s created", exercise.Name))
}

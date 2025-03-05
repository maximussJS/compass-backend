package services

import (
	"compass-backend/pkg/api/api_errors"
	dto_exercise "compass-backend/pkg/api/common/dto/exercise"
	"compass-backend/pkg/api/models"
	"compass-backend/pkg/api/repositories"
	fx_utils "compass-backend/pkg/common/fx"
	common_lib "compass-backend/pkg/common/lib"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type IExerciseService interface {
	Create(c *gin.Context, userId string, req dto_exercise.CreateExerciseRequest) (*models.Exercise, error)
}

type exerciseServiceParams struct {
	fx.In

	Logger                  common_lib.ILogger
	ExerciseRepository      repositories.IExerciseRepository
	ExerciseMediaRepository repositories.IExerciseMediaRepository
}

type exerciseService struct {
	logger                  common_lib.ILogger
	exerciseRepository      repositories.IExerciseRepository
	exerciseMediaRepository repositories.IExerciseMediaRepository
}

func FxExerciseService() fx.Option {
	return fx_utils.AsProvider(newExerciseService, new(IExerciseService))
}

func newExerciseService(params exerciseServiceParams) IExerciseService {
	return &exerciseService{
		logger:                  params.Logger,
		exerciseRepository:      params.ExerciseRepository,
		exerciseMediaRepository: params.ExerciseMediaRepository,
	}
}

func (s *exerciseService) Create(c *gin.Context, userId string, req dto_exercise.CreateExerciseRequest) (*models.Exercise, error) {
	existingExercise, existingErr := s.exerciseRepository.GetByNameAndCreatorId(c, req.Name, userId)

	if existingErr != nil {
		s.logger.Error(fmt.Sprintf("failed to get exercise by name and creator id: %v", existingErr))
		return nil, fmt.Errorf("failed to get exercise by name and creator id: %v", existingErr)
	}

	if existingExercise != nil {
		return nil, api_errors.ErrorExerciseAlreadyExists
	}

	exercise := models.Exercise{
		Name:        req.Name,
		Description: req.Description,
		CreatorId:   userId,
	}

	exerciseId, err := s.exerciseRepository.Create(c, exercise)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to create exercise: %v", err))
		return nil, fmt.Errorf("failed to create exercise: %v", err)
	}

	// Create exercise media

	newExercise, getErr := s.exerciseRepository.GetById(c, exerciseId)

	if getErr != nil {
		s.logger.Error(fmt.Sprintf("failed to get exercise by id: %v", getErr))
		return nil, fmt.Errorf("failed to get exercise by id: %v", getErr)
	}

	return newExercise, nil
}

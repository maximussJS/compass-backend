package services

import (
	"compass-backend/pkg/api/api_errors"
	category_dto "compass-backend/pkg/api/common/dto/category"
	"compass-backend/pkg/api/models"
	"compass-backend/pkg/api/repositories"
	fx_utils "compass-backend/pkg/common/fx"
	common_lib "compass-backend/pkg/common/lib"
	"context"
	"fmt"
	"go.uber.org/fx"
)

type ICategoryService interface {
	Create(ctx context.Context, dto category_dto.CreateCategoryRequest) (*models.Category, error)
	GetById(ctx context.Context, id uint) (*models.Category, error)
	List(ctx context.Context, limit, offset int) ([]models.Category, error)
	UpdateById(ctx context.Context, id uint, dto category_dto.UpdateCategoryRequest) (*models.Category, error)
	DeleteById(ctx context.Context, id uint) (*models.Category, error)
}

type categoryServiceParams struct {
	fx.In

	Logger             common_lib.ILogger
	CategoryRepository repositories.ICategoryRepository
}

type categoryService struct {
	logger common_lib.ILogger
	repo   repositories.ICategoryRepository
}

func FxCategoryService() fx.Option {
	return fx_utils.AsProvider(newCategoryService, new(ICategoryService))
}

func newCategoryService(params categoryServiceParams) *categoryService {
	return &categoryService{
		logger: params.Logger,
		repo:   params.CategoryRepository,
	}
}

func (s *categoryService) Create(ctx context.Context, dto category_dto.CreateCategoryRequest) (*models.Category, error) {
	existingCategory, existingErr := s.repo.GetByName(ctx, dto.Name)

	if existingErr != nil {
		s.logger.Error(fmt.Sprintf("failed to get category by name: %s", existingErr))
		return nil, existingErr
	}

	if existingCategory != nil {
		return existingCategory, api_errors.ErrorCategoryAlreadyExists
	}

	id, createErr := s.repo.Create(ctx, dto.ToModel())

	if createErr != nil {
		s.logger.Error(fmt.Sprintf("failed to create category: %s", createErr))
		return nil, createErr
	}

	category, err := s.repo.GetById(ctx, id)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to get category by id: %s", err))
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetById(ctx context.Context, id uint) (*models.Category, error) {
	category, err := s.repo.GetById(ctx, id)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to get category by id: %s", err))
		return nil, err
	}

	if category == nil {
		return nil, api_errors.ErrorCategoryNotFound
	}

	return category, nil
}

func (s *categoryService) List(ctx context.Context, limit, offset int) ([]models.Category, error) {
	categories, err := s.repo.List(ctx, limit, offset)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to list categories: %s", err))
		return nil, err
	}

	return categories, nil
}

func (s *categoryService) UpdateById(ctx context.Context, id uint, dto category_dto.UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.repo.GetById(ctx, id)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to get category by id: %s", err))
		return nil, err
	}

	if category == nil {
		return nil, api_errors.ErrorCategoryNotFound
	}

	if dto.Name != "" {
		existingCategory, existingErr := s.repo.GetByName(ctx, dto.Name)

		if existingErr != nil {
			s.logger.Error(fmt.Sprintf("failed to get category by name: %s", existingErr))
			return nil, existingErr
		}

		if existingCategory != nil {
			return existingCategory, api_errors.ErrorCategoryAlreadyExists
		}
	}

	updateErr := s.repo.UpdateById(ctx, id, dto.ToModel())

	if updateErr != nil {
		s.logger.Error(fmt.Sprintf("failed to update category: %s", updateErr))
		return nil, updateErr
	}

	updatedCategory, err := s.repo.GetById(ctx, id)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to get updated category by id: %s", err))
		return nil, err
	}

	return updatedCategory, nil
}

func (s *categoryService) DeleteById(ctx context.Context, id uint) (*models.Category, error) {
	category, err := s.repo.GetById(ctx, id)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to get category by id: %s", err))
		return nil, err
	}

	if category == nil {
		return nil, api_errors.ErrorCategoryNotFound
	}

	deleteErr := s.repo.DeleteById(ctx, id)

	if deleteErr != nil {
		s.logger.Error(fmt.Sprintf("failed to delete category by id: %s", deleteErr))
		return nil, deleteErr
	}

	return category, nil
}

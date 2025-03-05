package infrastructure

import (
	"compass-backend/pkg/api/models"
	common_infrastracture "compass-backend/pkg/common/infrastructure"
	common_lib "compass-backend/pkg/common/lib"
	common_models "compass-backend/pkg/common/models"
	"context"
	"fmt"
	"go.uber.org/fx"
)

var Module = fx.Options(
	FxCloudinary(),
	common_infrastracture.FxRouter(),
	common_infrastracture.FxRedis(),
	common_infrastracture.FxDatabase(),
	common_infrastracture.FxHttpServer(),
	fx.Invoke(func(
		lc fx.Lifecycle,
		logger common_lib.ILogger,
		database common_infrastracture.IDatabase,
	) {
		lc.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				instance := database.GetInstance()

				err := instance.SetupJoinTable(&common_models.Team{}, "Members", &common_models.TeamMember{})

				if err != nil {
					return fmt.Errorf("failed to run team members join table migration: %s", err)
				}
				err = instance.AutoMigrate(
					&common_models.User{},
					&common_models.Team{},
					&common_models.TeamMember{},
					&common_models.TeamInvite{},
					&models.Category{},
					&models.Exercise{},
					&models.ExerciseMedia{},
				)

				logger.Info("Migrations completed")
				return nil
			},
		})
	}),
)

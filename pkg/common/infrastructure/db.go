package infrastructure

import (
	"compass-backend/pkg/common/constants"
	fx_utils "compass-backend/pkg/common/fx"
	"compass-backend/pkg/common/lib"
	"context"
	"fmt"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type IDatabase interface {
	GetInstance() *gorm.DB
}

type databaseDependencies struct {
	fx.In

	Env    lib.IEnv
	Logger lib.ILogger
}

type database struct {
	logger   lib.ILogger
	env      lib.IEnv
	instance *gorm.DB
}

func FxDatabase() fx.Option {
	return fx_utils.AsProvider(newDatabase, new(IDatabase))
}

func newDatabase(lc fx.Lifecycle, deps databaseDependencies) IDatabase {
	db := &database{
		env:    deps.Env,
		logger: deps.Logger,
	}

	logMode := logger.Error

	if db.env.GetEnvironment() == constants.ProductionEnv {
		logMode = logger.Error
	}

	dbInstance, err := gorm.Open(postgres.Open(db.env.GetPostgresUrl()), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})

	if err != nil {
		db.logger.Error(fmt.Sprintf("Failed to connect to database: %s", err))
		panic(err)
	}

	db.logger.Info("Connected to database")

	db.instance = dbInstance

	lc.Append(fx.Hook{
		OnStop: db.shutdown,
	})

	return db
}

func (db *database) GetInstance() *gorm.DB {
	return db.instance
}

func (db *database) shutdown(_ context.Context) error {
	dbInstance, err := db.instance.DB()

	if err != nil {
		return fmt.Errorf("failed to get database instance: %s", err)
	}

	closeErr := dbInstance.Close()

	if closeErr != nil {
		return fmt.Errorf("failed to close database connection: %s", closeErr)
	}

	db.logger.Info("Database connection closed")

	return nil
}

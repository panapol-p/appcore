package appcore_store

import (
	"time"

	"github.com/panapol-p/appcore/appcore_utils"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDBStore(configs *appcore_utils.Configurations, logger *logrus.Logger) *gorm.DB {
	logger.Info("Connecting to database")
	db, err := gorm.Open(postgres.Open(configs.PostgresConnString), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Get generic database object sql.DB to use its functions
	sqlDB, _ := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	if configs.ObserveIsActive {
		if err = db.Use(otelgorm.NewPlugin()); err != nil {
			panic(err)
		}
	}

	logger.Info("Connecting to database success")
	return db
}

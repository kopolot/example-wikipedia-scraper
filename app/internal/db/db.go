package db

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	appLogger "example-wikipedia-scraper/internal/logger"

	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/pkg/db"
	"example-wikipedia-scraper/pkg/helpers"

	"github.com/golang-migrate/migrate/v4"
	postgresMig "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg config.ConfigInterface) error {
	dsn := cfg.GetDBString()
	logLevel := logger.Info
	if !cfg.IsDBDebugMode() {
		logLevel = logger.Error
	}
	loggerInstance := appLogger.GetLogger()
	writers := []io.Writer{}
	if loggerInstance != nil {
		writers = append(writers, loggerInstance.GetLogWriter())
	}
	if cfg.GetDBConfig().CliLogging {
		writers = append(writers, os.Stdout)
	}
	gormLogger := logger.New(
		log.New(io.MultiWriter(writers...), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("could not get DB instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return nil
}

func AutoMigrate() error {
	log.Println("Running database migrations...")
	sqlDB, _ := DB.DB()
	driver, err := postgresMig.WithInstance(sqlDB, &postgresMig.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", filepath.Join(helpers.GetCurrentFilePath(), "..", "..", "migrations")),
		"postgres",
		driver,
	)
	if err != nil {
		panic(err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	return nil
}

func RollbackAllMigrations() error {
	log.Println("Rolling back all database migrations...")
	sqlDB, _ := DB.DB()
	driver, err := postgresMig.WithInstance(sqlDB, &postgresMig.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", filepath.Join(helpers.GetCurrentFilePath(), "..", "..", "migrations")),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}
	if err = m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error during rollback: %w", err)
	}
	return nil
}

func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func GetQueryBuilder() db.QueryBuilder {
	return db.NewQueryBuilder(DB)
}

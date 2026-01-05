// Package persistence provides database implementations.
package persistence

import (
	"fmt"
	"log"

{{- if .EnableTelemetry}}
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
{{- end}}

	"{{.ModulePath}}/internal/infrastructure/config"
{{- if eq .DBDriver "postgres"}}
	"gorm.io/driver/postgres"
{{- else if eq .DBDriver "mysql"}}
	"gorm.io/driver/mysql"
{{- else if eq .DBDriver "sqlite"}}
	"gorm.io/driver/sqlite"
{{- end}}
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase creates a new GORM database connection
func NewDatabase(cfg config.DatabaseConfig) (*gorm.DB, error) {
	var dsn string
	var dialector gorm.Dialector

	switch cfg.Driver {
{{- if eq .DBDriver "postgres"}}
	case "postgres":
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Password,
			cfg.Name,
			cfg.SSLMode,
		)
		dialector = postgres.Open(dsn)
{{- else if eq .DBDriver "mysql"}}
	case "mysql":
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Name,
		)
		dialector = mysql.Open(dsn)
{{- else if eq .DBDriver "sqlite"}}
	case "sqlite", "sqlite3":
		dsn = cfg.Name
		dialector = sqlite.Open(dsn)
{{- end}}
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	// Configure GORM logger
	gormLogger := logger.Default.LogMode(logger.Silent)
	if cfg.Debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	// Open GORM connection
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
{{- if .EnableTelemetry}}

	// OpenTelemetry auto-instrumentation for database
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		return nil, fmt.Errorf("failed to enable otelgorm plugin: %w", err)
	}
{{- end}}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database connected successfully: %s:%s/%s", cfg.Host, cfg.Port, cfg.Name)

	return db, nil
}

// Transaction executes a function within a database transaction
func Transaction(db *gorm.DB, fn func(*gorm.DB) error) error {
	return db.Transaction(fn)
}

// AutoMigrate runs GORM auto migration for the given models
func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}

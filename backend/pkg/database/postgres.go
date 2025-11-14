package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/prabalesh/loco/backend/pkg/config"
	"go.uber.org/zap"
)

type Database struct {
	DB     *sql.DB
	Logger *zap.Logger
}

func NewPostgresDB(cfg config.DatabaseConfig, logger *zap.Logger) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	logger.Info("Connecting to database",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
		zap.String("database", cfg.Name),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("Failed to open database connection", zap.Error(err))
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxIdleConns(25)
	db.SetMaxIdleConns(5)

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connected successfully",
		zap.Int("max_open_conns", 25),
		zap.Int("max_idle_conns", 5),
		zap.Duration("conn_max_lifetime", 5*time.Minute),
	)

	return &Database{DB: db, Logger: logger}, nil
}

func (d *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := d.DB.QueryContext(ctx, query, args...)
	duration := time.Since(start)

	if duration > 500*time.Millisecond {
		d.Logger.Warn("Slow query detected",
			zap.Duration("duration", duration),
			zap.String("query", query),
		)
	}

	if err != nil {
		d.Logger.Error("Query failed",
			zap.Error(err),
			zap.String("query", query),
			zap.Duration("duration", duration),
		)
	}

	return rows, err
}

func (d *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := d.DB.ExecContext(ctx, query, args...)
	duration := time.Since(start)

	// Log slow queries
	if duration > 500*time.Millisecond {
		d.Logger.Warn("Slow query detected",
			zap.Duration("duration", duration),
			zap.String("query", query),
		)
	}

	if err != nil {
		d.Logger.Error("Exec failed",
			zap.Error(err),
			zap.String("query", query),
			zap.Duration("duration", duration),
		)
	}

	return result, err
}

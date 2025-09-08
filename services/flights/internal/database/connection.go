package database

import (
	"context"
	"fmt"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(databaseURL string) (*pgxpool.Pool, error) {

	if databaseURL == "" {
		logger.Error("database url is empty")
		return nil, fmt.Errorf("database url is empty")
	}

	config, err := pgxpool.ParseConfig(databaseURL)

	if err != nil {
		logger.Error("database url parse config error", "err", err)
		return nil, fmt.Errorf("error parsing DB config: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Error("Error creating the DB pool", "err", err)
		return nil, fmt.Errorf("error creating pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		logger.Error("Error pinging the DB pool", "err", err)
		pool.Close()
		return nil, fmt.Errorf("DB ping failed: %w", err)
	}

	logger.Info("Connected to DB")

	return pool, nil
}

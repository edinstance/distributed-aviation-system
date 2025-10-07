package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

func Init(cacheURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(cacheURL)

	if err != nil {
		logger.Error("Invalid Redis URL", "err", err)
		return nil, fmt.Errorf("invalid Redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//Tracing and Metrics
	if err := redisotel.InstrumentTracing(client); err != nil {
		return nil, fmt.Errorf("error setting up tracing: %w", err)
	}

	if err := redisotel.InstrumentMetrics(client); err != nil {
		return nil, fmt.Errorf("error setting up metrics: %w", err)
	}

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("Could not connect to Redis", "err", err)

		_ = client.Close()

		return nil, fmt.Errorf("could not connect to Redis: %w", err)
	}

	logger.Info("Connected to Redis")
	return client, nil
}

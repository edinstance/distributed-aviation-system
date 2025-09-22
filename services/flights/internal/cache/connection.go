package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
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

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("Could not connect to Redis", "err", err)
		return nil, fmt.Errorf("could not connect to Redis: %w", err)
	}

	logger.Info("Connected to Redis")
	return client, nil
}

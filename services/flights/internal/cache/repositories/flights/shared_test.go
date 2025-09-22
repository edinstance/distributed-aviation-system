package flights

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	val, ok := args.Get(0).(string)
	err := args.Error(1)
	if ok {
		return redis.NewStringResult(val, err)
	}
	return redis.NewStringResult("", err)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, exp time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, exp)
	return redis.NewStatusResult("", args.Error(0))
}

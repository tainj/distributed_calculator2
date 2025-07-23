package repository

import (
	"context"
	"fmt"
	"github.com/tainj/distributed_calculator2/pkg/db/cache"
)

type RedisResultRepository struct {
	cache *cache.CACHE
}

func NewRedisResultRepository(cache *cache.CACHE) *RedisResultRepository {
	return &RedisResultRepository{cache: cache}
}

func (r *RedisResultRepository) SetResult(ctx context.Context, variable string, result float64) error {
    err := r.cache.SetByKey(ctx, fmt.Sprintf("result:%s", variable), result)
	if err != nil {
		return fmt.Errorf("repository.SetResult: %w", err)
	}
	return nil
}

func (r *RedisResultRepository) GetResult(ctx context.Context, variable string, result float64) error {
    err := r.cache.GetByKey(ctx, fmt.Sprintf("result:%s", variable), result)
	if err != nil {
		return fmt.Errorf("repository.GetResult: %w", err)
	}
	return nil
}

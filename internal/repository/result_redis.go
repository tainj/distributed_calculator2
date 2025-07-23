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

func (r *RedisResultRepository) SaveResult(ctx context.Context, variable string, result float64) error {
    return r.cache.SetByKey(ctx, fmt.Sprintf("result:%s", variable), result)
}
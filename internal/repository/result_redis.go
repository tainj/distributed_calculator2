package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/tainj/distributed_calculator2/internal/models"
	"github.com/tainj/distributed_calculator2/pkg/db/cache"
)

type RedisResultRepository struct {
	cache *cache.CACHE
}

func NewRedisResultRepository(cache *cache.CACHE) *RedisResultRepository {
	return &RedisResultRepository{cache: cache}
}


func (s *RedisResultRepository) GetExpressions(ctx context.Context) ([]models.Example, error) {
	// получаем данные из Redis
	key := "user:1:expressions"
	data, err := s.cache.Client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Ключ не найден - не ошибка
		}
		return nil, fmt.Errorf("failed to get from Redis: %w", err)
	}

	// десериализуем из JSON
	var expressions []models.Example
	if err := json.Unmarshal(data, &expressions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal expressions: %w", err)
	}

	return expressions, nil
}

func (s *RedisResultRepository) AddExpression(ctx context.Context, expr models.Example) error {
	// получаем существующие
	key := "user:1:expressions"
	expressions, err := s.GetExpressions(ctx)
	if err != nil {
		return err
	}

	// добавляем новое
	expressions = append(expressions, expr)

	// сохраняем обновленный список
	return s.SaveExpressions(ctx, key, expressions)
}

func (s *RedisResultRepository) SaveExpressions(ctx context.Context, key string, expressions []models.Example) error {
	// сериализуем в JSON
	data, err := json.Marshal(expressions)
	if err != nil {
		return fmt.Errorf("failed to marshal expressions: %w", err)
	}

	// записываем в Redis
	if err := s.cache.Client.Set(ctx, key, data, 0).Err(); err != nil {
		return fmt.Errorf("failed to save to Redis: %w", err)
	}

	return nil
}

func (r *RedisResultRepository) SaveResult(ctx context.Context, variable string, result float64) error {
    return r.cache.SetByKey(ctx, fmt.Sprintf("user:1:variable:%s", variable), result)
}
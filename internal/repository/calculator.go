package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/tainj/distributed_calculator2/internal/models"
	"github.com/tainj/distributed_calculator2/pkg/db/cache"
	"github.com/tainj/distributed_calculator2/pkg/db/postgres"
)

type CalculatorRepository struct {
	db *postgres.DB
	cache *cache.CACHE
}

func NewCalculatorRepository(db *postgres.DB, cache *cache.CACHE) *CalculatorRepository {
	return &CalculatorRepository{db: db, cache: cache}
}

// func (r *CalculatorRepository) SaveResult(expression string, result float64) error {
// 	log.Printf("Saved result: %s = %v", expression, result)
// 	return nil
// }

// func (r *CalculatorRepository) GetResult(expression string) (float64, error) {
// 	return 42, nil
// }

func (r *CalculatorRepository) SaveExamples() {

}

func (r *CalculatorRepository) GetExpressions(ctx context.Context) ([]models.Example, error) {
	// получаем данные из Redis
	key := "user:1:expressions"
	data, err := r.cache.Client.Get(ctx, key).Bytes()
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

func (r *CalculatorRepository) AddExpression(ctx context.Context, expr models.Example) error {
	// получаем существующие
	key := "user:1:expressions"
	expressions, err := r.GetExpressions(ctx)
	if err != nil {
		return err
	}

	// добавляем новое
	expressions = append(expressions, expr)

	// сохраняем обновленный список
	return r.SaveExpressions(ctx, key, expressions)
}

func (r *CalculatorRepository) SaveExpressions(ctx context.Context, key string, expressions []models.Example) error {
	// сериализуем в JSON
	data, err := json.Marshal(expressions)
	if err != nil {
		return fmt.Errorf("failed to marshal expressions: %w", err)
	}

	// записываем в Redis
	if err := r.cache.Client.Set(ctx, key, data, 0).Err(); err != nil {
		return fmt.Errorf("failed to save to Redis: %w", err)
	}

	return nil
}
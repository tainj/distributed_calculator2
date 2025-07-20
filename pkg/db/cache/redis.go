package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/tainj/distributed_calculator2/internal/models"
)

type Config struct {
	Host string `env:"REDIS_HOST" env-default:"localhost"`
	Port string `env:"REDIS_PORT" env-default:"6379"`
}

type CACHE struct {
	Client *redis.Client
}

func New(cfg Config) *CACHE {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	})

	// Простая проверка без таймаута
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Printf("Redis connection warning: %v", err) // не фатальная ошибка
	} else {
		log.Println("Redis connected")
	}

	return &CACHE{Client: client}
}

func (s *CACHE) GetExpressions(ctx context.Context) ([]models.Example, error) {
	// получаем данные из Redis
	key := "user:1:expressions"
	data, err := s.Client.Get(ctx, key).Bytes()
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

func (s *CACHE) AddExpression(ctx context.Context, expr models.Example) error {
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

func (s *CACHE) SaveExpressions(ctx context.Context, key string, expressions []models.Example) error {
	// сериализуем в JSON
	data, err := json.Marshal(expressions)
	if err != nil {
		return fmt.Errorf("failed to marshal expressions: %w", err)
	}

	// записываем в Redis
	if err := s.Client.Set(ctx, key, data, 0).Err(); err != nil {
		return fmt.Errorf("failed to save to Redis: %w", err)
	}

	return nil
}


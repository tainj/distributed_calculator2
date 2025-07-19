package cache

import (
	"context"
	"fmt"
	"log"
	"github.com/redis/go-redis/v9"
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


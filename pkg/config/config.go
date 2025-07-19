package config

import (
	"github.com/tainj/distributed_calculator2/pkg/db/cache"
	"github.com/tainj/distributed_calculator2/pkg/db/postgres"
)

type Config struct {
	postgres.Config
	cache.RedisConfig
}


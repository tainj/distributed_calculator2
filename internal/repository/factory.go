package repository

import (
	"github.com/tainj/distributed_calculator2/pkg/db/cache"
	"github.com/tainj/distributed_calculator2/pkg/db/postgres"
    postgresRepo "github.com/tainj/distributed_calculator2/internal/repository/postgres"
    redisRepo "github.com/tainj/distributed_calculator2/internal/repository/redis"
)


type RepositoryFactory struct {
    postgresDB *postgres.DB
    redisCache *cache.CACHE
}

func NewRepositoryFactory(postgresDB *postgres.DB, redisCache *cache.CACHE) *RepositoryFactory {
    return &RepositoryFactory{
        postgresDB: postgresDB,
        redisCache: redisCache,
    }
}

// CreateUserRepository создает репозиторий для работы с пользователями
func (f *RepositoryFactory) CreateUserRepository() UserRepository {
    return postgresRepo.NewAuthUserRepository(f.postgresDB)
}

// CreateExampleRepository создает репозиторий для работы с примерами
func (f *RepositoryFactory) CreateExampleRepository() ExampleRepository {
    return postgresRepo.NewPostgresResultRepository(f.postgresDB)
}

// CreateVariableRepository создает репозиторий для работы с переменными в Redis
func (f *RepositoryFactory) CreateVariableRepository() VariableRepository {
    return redisRepo.NewRedisResultRepository(f.redisCache)
}
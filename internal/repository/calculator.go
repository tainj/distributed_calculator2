package repository

import (
	"github.com/tainj/distributed_calculator2/pkg/db/cache"
	"github.com/tainj/distributed_calculator2/pkg/db/postgres"
)

type CalculatorRepository struct {
	Db *postgres.DB
	Cache *cache.CACHE
}

func NewCalculatorRepository(db *postgres.DB, cache *cache.CACHE) *CalculatorRepository {
	return &CalculatorRepository{Db: db, Cache: cache}
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
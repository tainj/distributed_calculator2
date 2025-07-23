package repository

import (
	"context"

	"github.com/tainj/distributed_calculator2/internal/models"
	"github.com/tainj/distributed_calculator2/pkg/db/postgres"
)

type PostgresResultRepository struct {
	db *postgres.DB
}

func NewPostgresResultRepository(db *postgres.DB) *PostgresResultRepository {
	return &PostgresResultRepository{db: db}
}

// func (s *PostgresResultRepository)SaveExample(ctx context.Context, example *models.Example) error {

// }
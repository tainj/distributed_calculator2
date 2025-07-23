package repository

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/tainj/distributed_calculator2/internal/models"
	"github.com/tainj/distributed_calculator2/pkg/db/postgres"
)

type PostgresResultRepository struct {
	db *postgres.DB
}

func NewPostgresResultRepository(db *postgres.DB) *PostgresResultRepository {
	return &PostgresResultRepository{db: db}
}

func (r *PostgresResultRepository) SaveExample(ctx context.Context, example *models.Example) error {
    query := sq.Insert("examples").
        Columns("id", "expression", "response", "user_id", "calculated").
        Values(example.Id, example.Expression, example.Response, "1", false).
        Suffix("RETURNING id").
        PlaceholderFormat(sq.Dollar).
        RunWith(r.db.Db)

    var insertedID string
    err := query.QueryRowContext(ctx).Scan(&insertedID)
    if err != nil {
        return fmt.Errorf("repository.SaveExample: failed to insert example: %w", err)
    }

    return nil
}
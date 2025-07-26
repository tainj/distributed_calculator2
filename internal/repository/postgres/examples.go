package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/tainj/distributed_calculator2/internal/models"
	"github.com/tainj/distributed_calculator2/pkg/db/postgres"
)

func NewPostgresResultRepository(db *postgres.DB) *PostgresResultRepository {
    return &PostgresResultRepository{db: db}
}

type PostgresResultRepository struct {
	db *postgres.DB
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

func (r *PostgresResultRepository) UpdateExample(ctx context.Context, exampleId string, result float64) error {
    query := sq.Update("examples").
        Set("calculated", true).
        Set("result", result).
        Where(sq.Eq{"id": exampleId}).
        PlaceholderFormat(sq.Dollar).
        RunWith(r.db.Db)

    _, err := query.ExecContext(ctx)
    if err != nil {
        return fmt.Errorf("repository.UpdateExample: %w", err)
    }

    return nil
}

func (r *PostgresResultRepository) GetResult(ctx context.Context, exampleID string) (float64, error) {
    query := sq.Select("result").
        From("examples").
        Where(sq.Eq{"id": exampleID}).
        PlaceholderFormat(sq.Dollar).
        RunWith(r.db.Db)

    var example models.UserExample
    err := query.QueryRowContext(ctx).Scan(
        &example.Result,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return 0, fmt.Errorf("repository.GetResult: example not found")
        }
        return 0, fmt.Errorf("repository.GetResult: failed to get example: %w", err)
    }

    return *example.Result, nil
}
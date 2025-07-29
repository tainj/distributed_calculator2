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
    calculated := example.Error != nil

    query := sq.Insert("examples").
        Columns("id", "expression", "response", "user_id", "calculated", "error").
        Values(
            example.ID,
            example.Expression,
            example.Response,
            example.UserID,
            calculated,
            example.Error,
        ).
        PlaceholderFormat(sq.Dollar).
        RunWith(r.db.Db)

    _, err := query.ExecContext(ctx)
    if err != nil {
        return fmt.Errorf("repository.SaveExample: failed to insert example: %w", err)
    }
    
    return nil
}

func (r *PostgresResultRepository) UpdateExample(ctx context.Context, exampleId string, result float64) error {
    // обновляем статус и результат по id
    query := sq.Update("examples").
        Set("calculated", true).
        Set("result", result).
        Where(sq.Eq{"id": exampleId}).
        PlaceholderFormat(sq.Dollar).
        RunWith(r.db.Db)

    // применяем изменения
    _, err := query.ExecContext(ctx)
    if err != nil {
        return fmt.Errorf("repository.UpdateExample: %w", err)
    }
    
    return nil
}

func (r *PostgresResultRepository) UpdateExampleWithError(ctx context.Context, exampleID, errorMsg string) error {
    query := sq.Update("examples").
        Set("calculated", true).
        Set("error", errorMsg).
        Where(sq.Eq{"id": exampleID}).
        PlaceholderFormat(sq.Dollar).
        RunWith(r.db.Db)

    _, err := query.ExecContext(ctx)
    if err != nil {
        return fmt.Errorf("failed to update example with error: %w", err)
    }
    return nil
}

func (r *PostgresResultRepository) GetResult(ctx context.Context, exampleID string) (float64, error) {
    var calculated bool
    var result sql.NullFloat64
    var dbError sql.NullString

    query := sq.Select("calculated", "result", "error").
        From("examples").
        Where(sq.Eq{"id": exampleID}).
        PlaceholderFormat(sq.Dollar).
        RunWith(r.db.Db)

    err := query.QueryRowContext(ctx).Scan(&calculated, &result, &dbError)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return 0, fmt.Errorf("example not found")
        }
        return 0, fmt.Errorf("failed to query: %w", err)
    }

    if !calculated {
        return 0, fmt.Errorf("calculation not completed yet")
    }

    if dbError.Valid {
        return 0, fmt.Errorf("calculation failed: %s", dbError.String)
    }

    if !result.Valid {
        return 0, fmt.Errorf("result is not available")
    }

    return result.Float64, nil
}

func (r *PostgresResultRepository) GetExamplesByUserID(ctx context.Context, userID string) ([]models.Example, error) {
    // строим запрос через squirrel
    query := sq.Select("id", "expression", "calculated", "result", "error", "created_at").
        From("examples").
        Where(sq.Eq{"user_id": userID}).
        OrderBy("created_at DESC").
        PlaceholderFormat(sq.Dollar)

    // выполняем
    rows, err := query.RunWith(r.db.Db).QueryContext(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to query examples: %w", err)
    }
    defer rows.Close()

    var examples []models.Example
    for rows.Next() {
        var example models.Example
        var result sql.NullFloat64
        var dbError sql.NullString

        err := rows.Scan(
            &example.ID,
            &example.Expression,
            &example.Calculated,
            &result,
            &dbError,
            &example.CreatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan row: %w", err)
        }

        // если есть результат — сохраняем
        if result.Valid {
            example.Result = &result.Float64
        }

        // если есть ошибка — сохраняем
        if dbError.Valid {
            example.Error = &dbError.String
        }

        examples = append(examples, example)
    }

    // проверяем, не было ли ошибки в rows
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("row iteration error: %w", err)
    }

    return examples, nil
}
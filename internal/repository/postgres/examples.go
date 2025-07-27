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
    // формируем запрос для сохранения примера в бд
    query := sq.Insert("examples").
        Columns("id", "expression", "response", "user_id", "calculated").
        Values(example.Id, example.Expression, example.Response, example.UserID, false).
        PlaceholderFormat(sq.Dollar).
        RunWith(r.db.Db)

    // выполняем вставку, проверяем ошибки
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
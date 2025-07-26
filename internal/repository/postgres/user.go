package repository

import (
    "context"
    "fmt"

    sq "github.com/Masterminds/squirrel"
    "github.com/tainj/distributed_calculator2/internal/models"
    "github.com/tainj/distributed_calculator2/pkg/db/postgres"
)

type AuthUserRepository struct {
    db *postgres.DB
}

func NewAuthUserRepository(db *postgres.DB) *AuthUserRepository {
    return &AuthUserRepository{db: db}
}

func (r *AuthUserRepository) Register(ctx context.Context, user *models.User) error {
    query := sq.Insert("users").
        Columns("id", "username", "email", "password_hash", "role").
        Values(user.ID, user.Username, user.Email, user.PasswordHash, string(user.Role)).
        PlaceholderFormat(sq.Dollar)

    sql, args, err := query.ToSql()
    if err != nil {
        return fmt.Errorf("failed to build query: %w", err)
    }

    _, err = r.db.Db.ExecContext(ctx, sql, args...)
    if err != nil {
        return fmt.Errorf("failed to insert user: %w", err)
    }

    return nil
}

func (r *AuthUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
    query := sq.Select("id", "username", "email", "password_hash", "role", "created_at", "updated_at").
        From("users").
        Where(sq.Eq{"username": username}).
        PlaceholderFormat(sq.Dollar)

    sql, args, err := query.ToSql()
    if err != nil {
        return nil, fmt.Errorf("failed to build query: %w", err)
    }

    var user models.User
    err = r.db.Db.QueryRowContext(ctx, sql, args...).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.Role,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to get user by username: %w", err)
    }

    return &user, nil
}

func (r *AuthUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    query := sq.Select("id", "username", "email", "password_hash", "role", "created_at", "updated_at").
        From("users").
        Where(sq.Eq{"email": email}).
        PlaceholderFormat(sq.Dollar)

    sql, args, err := query.ToSql()
    if err != nil {
        return nil, fmt.Errorf("failed to build query: %w", err)
    }

    var user models.User
    err = r.db.Db.QueryRowContext(ctx, sql, args...).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.Role,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to get user by email: %w", err)
    }

    return &user, nil
}

func (r *AuthUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
    query := sq.Select("id", "username", "email", "password_hash", "role", "created_at", "updated_at").
        From("users").
        Where(sq.Eq{"id": id}).
        PlaceholderFormat(sq.Dollar)

    sql, args, err := query.ToSql()
    if err != nil {
        return nil, fmt.Errorf("failed to build query: %w", err)
    }

    var user models.User
    err = r.db.Db.QueryRowContext(ctx, sql, args...).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.Role,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to get user by id: %w", err)
    }

    return &user, nil
}
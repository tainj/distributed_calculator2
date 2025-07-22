package repository

import (
	"context"

	"github.com/tainj/distributed_calculator2/internal/models"
)

type ResultRepository interface {
	GetExpressions(ctx context.Context) ([]models.Example, error)
	AddExpression(ctx context.Context, expr models.Example) error
	SaveExpressions(ctx context.Context, key string, expressions []models.Example) error
	SaveResult(ctx context.Context, variable string, result float64) error
}
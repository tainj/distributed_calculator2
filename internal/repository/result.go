package repository

import (
	"context"
)

type ResultRepository interface {
	SaveResult(ctx context.Context, variable string, result float64) error
}
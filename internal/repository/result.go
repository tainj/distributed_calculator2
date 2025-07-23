package repository

import (
	"context"

	"github.com/tainj/distributed_calculator2/internal/models"
)

type ResultRepository interface {
	SaveResult(ctx context.Context, variable string, result float64) error
}

type ExampleRepository interface {
	SaveExample(ctx context.Context, example *models.Example) error
}
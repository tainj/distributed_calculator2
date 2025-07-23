package repository

import (
	"context"

	"github.com/tainj/distributed_calculator2/internal/models"
)

type CacheRepository interface {
    SetResult(ctx context.Context, variable string, result float64) error
    // GetResult(ctx context.Context, variable string, dest *float64) error
}

type ExampleRepository interface {
    SaveExample(ctx context.Context, example *models.Example) error
    UpdateExample(ctx context.Context, exampleId string, result float64) error
    GetResult(ctx context.Context, exampleID string) (float64, error)
    // GetExamples(ctx context.Context, userID string) ([]models.Example, error)
}
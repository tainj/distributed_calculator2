package grpc

import (
	"context"
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/tainj/distributed_calculator2/internal/models"
	client "github.com/tainj/distributed_calculator2/pkg/api"
)

type Service interface {
	Calculate(ctx context.Context, example *models.Example) (*models.Example, error)
	GetResult(ctx context.Context, exampleID string) (float64, error)
}

type CalculatorService struct {
	client.UnimplementedCalculatorServer
	service Service
}

func NewOrderService(srv Service) *CalculatorService {
	return &CalculatorService{service: srv}
}

func (s *CalculatorService) Calculate(ctx context.Context, req *client.CalculateRequest) (*client.CalculateResponse, error) {
	resp, err := s.service.Calculate(ctx, &models.Example{
		Expression: req.GetExpression(),
	})
	if err != nil {
		return nil, fmt.Errorf("Calculate: %w", err)
	}
	r := pointer.Get(resp)
	return &client.CalculateResponse{
		TaskId: r.Id,
	}, nil
}

// transport/grpc/calculator.go
func (s *CalculatorService) GetResult(ctx context.Context, req *client.GetResultRequest) (*client.GetResultResponse, error) {
    taskID := req.GetTaskId()

    // Получаем из бизнес-логики
    result, err := s.service.GetResult(ctx, taskID)
    if err != nil {
        return &client.GetResultResponse{
            Result: &client.GetResultResponse_Error{
                Error: err.Error(),
            },
        }, nil // ❗ не возвращаем ошибку в gRPC sense, а кладём в oneof
    }

    // Успешно — возвращаем значение
    return &client.GetResultResponse{
        Result: &client.GetResultResponse_Value{
            Value: result, // *float64 → float64
        },
    }, nil
}
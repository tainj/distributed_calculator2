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
}

type CalculatorService struct {
	client.UnimplementedCalculatorServer
	service Service
}

func NewOrderService(srv Service) *CalculatorService {
	return &CalculatorService{service: srv}
}

func (s *CalculatorService) Calculate(ctx context.Context, req *client.CalculationRequest) (*client.CalculationResponse, error) {
	resp, err := s.service.Calculate(ctx, &models.Example{
		Expression: req.GetExpression(),
	})
	if err != nil {
		return nil, fmt.Errorf("Calculate: %w", err)
	}
	r := pointer.Get(resp)
	return &client.CalculationResponse{
		TaskId: r.Response,
	}, nil
}

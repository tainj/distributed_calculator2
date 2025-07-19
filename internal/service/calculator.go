package service

import (
	"context"
	"log"
	repo "github.com/tainj/distributed_calculator2/internal/repository"
)

type CalculatorService struct {
	repo repo.CalculatorRepository
}

func NewCalculatorService(repo repo.CalculatorRepository) *CalculatorService {
	return &CalculatorService{repo: repo}
}

func (s *CalculatorService) Calculate(ctx context.Context, req *api.CalculationRequest) (*api.CalculationResponse, error) {
	// Заглушка: просто возвращаем результат
	log.Printf("Calculating expression: %s", req.Expression)
	result := 42.0 // В реальности здесь будет вычисление

	// Сохраняем результат в репозиторий
	if err := s.repo.SaveResult(req.Expression, result); err != nil {
	return nil, err
	}

	return &api.CalculationResponse{
	Result: &api.CalculationResponse_Value{Value: result},
	}, nil
}
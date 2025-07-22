package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tainj/distributed_calculator2/internal/models"
	repo "github.com/tainj/distributed_calculator2/internal/repository"
	"github.com/tainj/distributed_calculator2/kafka"
	"github.com/tainj/distributed_calculator2/pkg/calculator"
)

type CalculatorService struct {
	repo *repo.CalculatorRepository
	kafkaQueue  *kafka.KafkaQueue
}

func NewCalculatorService(repo *repo.CalculatorRepository, kafkaQueue *kafka.KafkaQueue) *CalculatorService {
	return &CalculatorService{repo: repo, kafkaQueue: kafkaQueue}
}

func (s *CalculatorService) Calculate(ctx context.Context, example *models.Example) (*models.Example, error) {
	expr := calculator.NewExample(example.Expression)  // создаем структуру
	if !expr.Check() {
		return nil, models.ErrCovertExample
	}
	expr.Convert() // переводим в польскую нотацию
	results, variable := expr.Calculate()
	example = &models.Example{
		Id: uuid.NewString(),
		Expression: example.Expression,
		SimpleExamples: results,
		Response: variable,
	}

	for _, task := range results {
		err := s.kafkaQueue.SendTask(task)
		return nil, fmt.Errorf("send message kafka: %w", err)
	}

	s.repo.Cache.AddExpression(ctx, *example)
	return example, nil
}
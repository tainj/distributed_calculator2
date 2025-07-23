package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tainj/distributed_calculator2/internal/models"
	repo "github.com/tainj/distributed_calculator2/internal/repository"
	"github.com/tainj/distributed_calculator2/pkg/messaging/kafka"
	"github.com/tainj/distributed_calculator2/pkg/calculator"
)

type CalculatorService struct {
	repo repo.RedisResultRepository
	repoExamples repo.ExampleRepository
	kafkaQueue kafka.TaskQueue
}

func NewCalculatorService(repo *repo.RedisResultRepository, kafkaQueue *kafka.KafkaQueue, repoExample *repo.PostgresResultRepository) *CalculatorService {
	return &CalculatorService{repo: *repo, kafkaQueue: kafkaQueue, repoExamples: repoExample}
}

func (s *CalculatorService) Calculate(ctx context.Context, example *models.Example) (*models.Example, error) {
	expr := calculator.NewExpression(example.Expression)  // создаем структуру
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

	err := s.repoExamples.SaveExample(ctx, example)
	if err != nil {
		return nil, fmt.Errorf("Calculate: save example: %v", err)
	}

	for _, task := range results {
		if err := s.kafkaQueue.SendTask(task); err != nil {
        	return nil, fmt.Errorf("Calculate: send message kafka: %v", err)
    	}
	}

	return example, nil
}
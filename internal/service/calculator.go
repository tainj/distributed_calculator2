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
	repoExamples repo.ExampleRepository
	kafkaQueue kafka.TaskQueue
}

func NewCalculatorService(kafkaQueue kafka.TaskQueue, repoExample repo.ExampleRepository) *CalculatorService {
	return &CalculatorService{kafkaQueue: kafkaQueue, repoExamples: repoExample}
}

func (s *CalculatorService) Calculate(ctx context.Context, example *models.Example) (*models.Example, error) {
	expr := calculator.NewExpression(example.Expression)  // создаем структуру
	if !expr.Check() {
		return nil, models.ErrCovertExample
	}
	if _, err := expr.Convert(); err != nil {  // переводим в польскую нотацию 
        return nil, err
    }

	// генерируем ID примера
    exampleID := uuid.New().String()

	results, variable := expr.Calculate()

	example = &models.Example{ // формируем пример для сохранения
		Id: exampleID,
		Expression: example.Expression,
		SimpleExamples: results,
		Response: variable,
	}

	err := s.repoExamples.SaveExample(ctx, example) // сохраняем
	if err != nil {
		return nil, fmt.Errorf("Calculate: save example: %v", err)
	}

	for i, task := range results { // отправляем по очереди таски
		kafkaTask := &models.Task{
            Num1:      task.Num1,
            Num2:      task.Num2,
            Sign:      task.Sign,
            Variable:  task.Variable,
            ExampleID: exampleID,
            Index:     i,
            IsFinal:   task.Variable == variable,
        }

        if err := s.kafkaQueue.SendTask(kafkaTask); err != nil {
            return nil, fmt.Errorf("failed to send task to kafka: %w", err)
        }
	}
	return example, nil
}

func (s *CalculatorService) GetResult(ctx context.Context, exampleID string) (float64, error) {
	return s.repoExamples.GetResult(ctx, exampleID)
}
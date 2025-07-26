package service

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/tainj/distributed_calculator2/internal/auth"
    "github.com/tainj/distributed_calculator2/internal/models"
    repo "github.com/tainj/distributed_calculator2/internal/repository"
    "github.com/tainj/distributed_calculator2/pkg/calculator"
    "github.com/tainj/distributed_calculator2/pkg/messaging/kafka"
)

// calculator service — orchestrator
// отправляет задачи в кафку, сохраняет примеры
type CalculatorService struct {
    userRepo     repo.UserRepository
    repoExamples repo.ExampleRepository
    kafkaQueue   kafka.TaskQueue
}

// newcalculator service
func NewCalculatorService(userRepo repo.UserRepository, kafkaQueue kafka.TaskQueue, repoExample repo.ExampleRepository) *CalculatorService {
    return &CalculatorService{
        kafkaQueue:   kafkaQueue,
        repoExamples: repoExample,
		userRepo: userRepo,
    }
}

// calculate — запускает вычисление выражения
func (s *CalculatorService) Calculate(ctx context.Context, example *models.Example) (*models.Example, error) {
    // создаём парсер выражения
    expr := calculator.NewExpression(example.Expression)
    if !expr.Check() {
        return nil, models.ErrCovertExample
    }

    // переводим в польскую нотацию
    if _, err := expr.Convert(); err != nil {
        return nil, err
    }

    // генерируем id
    exampleID := uuid.New().String()

    // считаем шаги и финальную переменную
    results, variable := expr.Calculate()

    // формируем пример
    example = &models.Example{
        Id:             exampleID,
        Expression:     example.Expression,
        SimpleExamples: results,
        Response:       variable,
    }

    // сохраняем в бд
    err := s.repoExamples.SaveExample(ctx, example)
    if err != nil {
        return nil, fmt.Errorf("calculate: save example: %v", err)
    }

    // отправляем каждый шаг в кафку
    for i, task := range results {
        kafkaTask := &models.Task{
            Num1:      task.Num1,
            Num2:      task.Num2,
            Sign:      task.Sign,
            Variable:  task.Variable,
            ExampleID: exampleID,
            Index:     i,
            IsFinal:   task.Variable == variable, // последний шаг
        }

        if err := s.kafkaQueue.SendTask(kafkaTask); err != nil {
            return nil, fmt.Errorf("failed to send task to kafka: %w", err)
        }
    }

    return example, nil
}

// getresult — получает финальный результат по id
func (s *CalculatorService) GetResult(ctx context.Context, exampleID string) (float64, error) {
    return s.repoExamples.GetResult(ctx, exampleID)
}

// register — регистрирует нового пользователя
func (s *CalculatorService) Register(ctx context.Context, user *models.UserCredentials) (*models.User, error) {
    // проверяем, есть ли уже с такой почтой
    if _, err := s.userRepo.GetByEmail(ctx, user.Email); err == nil {
        return nil, fmt.Errorf("email already exists")
    }

    // хэшируем пароль
    hashedPassword, err := auth.HashPassword(user.Password)
    if err != nil {
        return nil, err
    }

    // создаём нового пользователя
    newUser := &models.User{
        ID:           uuid.New().String(),
        Email:        user.Email,
        PasswordHash: hashedPassword,
        Role:         models.UserRole,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    // сохраняем
    if err := s.userRepo.Register(ctx, newUser); err != nil {
        return nil, err
    }

    return newUser, nil
}
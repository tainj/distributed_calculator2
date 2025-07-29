package service

import (
	"context"
	"errors"
	"fmt"
	"log"
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
    jwtService   auth.JWTService
}

// newcalculator service
func NewCalculatorService(userRepo repo.UserRepository, kafkaQueue kafka.TaskQueue, repoExample repo.ExampleRepository, jwtService   auth.JWTService) *CalculatorService {
    return &CalculatorService{
        kafkaQueue:   kafkaQueue,
        repoExamples: repoExample,
		userRepo: userRepo,
        jwtService: jwtService,
    }
}

// calculate — запускает вычисление выражения
func (s *CalculatorService) Calculate(ctx context.Context, example *models.Example) (*models.Example, error) {
    exampleID := uuid.New().String()

    resultExample := &models.Example{
        ID:         exampleID,
        Expression: example.Expression,
        UserID:     example.UserID,
    }

    // создаём парсер выражения
    expr := calculator.NewExpression(example.Expression)

    // переводим в польскую нотацию
    if _, err := expr.Convert(); err != nil {
        errString := err.Error()
        resultExample.Error = &errString

        // 🔽 ДОБАВЬ ЭТО
        log.Printf("[DEBUG] Saving example with error: ID=%s, Expr=%s, Error=%s, Response=%q", 
            resultExample.ID, resultExample.Expression, errString, resultExample.Response)

        if errSave := s.repoExamples.SaveExample(ctx, resultExample); errSave != nil {
            log.Printf("[ERROR] Failed to save example with error: %v", errSave)
            return nil, fmt.Errorf("calculate: save example: %v", errSave)
        }
        return resultExample, err
    }

    // считаем шаги и финальную переменную
    results, variable := expr.Calculate()

    // заполняем результат
    resultExample.SimpleExamples = results
    resultExample.Response = variable // не хватало этого!

    if err := s.repoExamples.SaveExample(ctx, resultExample); err != nil {
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
            IsFinal:   task.Variable == variable,
        }

        if err := s.kafkaQueue.SendTask(kafkaTask); err != nil {
            return nil, fmt.Errorf("failed to send task to kafka: %w", err)
        }
    }

    return resultExample, nil
}

// getresult — получает финальный результат по id
func (s *CalculatorService) GetResult(ctx context.Context, exampleID string) (float64, error) {
    return s.repoExamples.GetResult(ctx, exampleID)
}

// register — регистрирует нового пользователя
func (s *CalculatorService) Register(ctx context.Context, userRequest *models.UserCredentials) (*models.User, error) {
    // проверяем, есть ли уже с такой почтой
    if _, err := s.userRepo.GetByEmail(ctx, userRequest.Email); err == nil {
        return nil, fmt.Errorf("email already exists")
    }

    // хэшируем пароль
    hashedPassword, err := auth.HashPassword(userRequest.Password)
    if err != nil {
        return nil, err
    }

    // создаём нового пользователя
    user := &models.User{
        ID:           uuid.New().String(),
        Email:        userRequest.Email,
        PasswordHash: hashedPassword,
        Role:         models.UserRole,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    // сохраняем
    if err := s.userRepo.Register(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

func (s *CalculatorService) Login(ctx context.Context, userRequest *models.UserCredentials) (*models.LoginResponse, error) {
    user, err := s.userRepo.GetByEmail(ctx, userRequest.Email)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    if !auth.CheckPassword(userRequest.Password, user.PasswordHash) {
        return nil, errors.New("invalid credentials")
    }

    // генерируем JWT
    token, err := s.jwtService.GenerateToken(user.ID)
    if err != nil {
        return nil, err
    }

    // возвращаем ответ
    return &models.LoginResponse{
        UserID: user.ID,
        Email:  user.Email,
        Token:  token,
    }, nil
}

func (s *CalculatorService) GetExamplesByUserID(ctx context.Context, userID string) ([]models.Example, error) {
    return s.repoExamples.GetExamplesByUserID(ctx, userID)
}
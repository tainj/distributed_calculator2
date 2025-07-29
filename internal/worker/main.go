// internal/worker/worker.go
package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tainj/distributed_calculator2/internal/models"
	repo "github.com/tainj/distributed_calculator2/internal/repository"
	"github.com/tainj/distributed_calculator2/internal/valueprovider"
	"github.com/tainj/distributed_calculator2/pkg/calculator"
	"github.com/tainj/distributed_calculator2/pkg/logger"
	"github.com/tainj/distributed_calculator2/pkg/messaging/kafka"
)

type Worker struct {
    exampleRepo   repo.ExampleRepository
    cacheRepo     repo.VariableRepository
    kafkaQueue    kafka.TaskQueue
    valueProvider valueprovider.Provider
    logger        logger.Logger
}

func NewWorker(
    exampleRepo repo.ExampleRepository,
    cacheRepo repo.VariableRepository,
    kafkaQueue kafka.TaskQueue,
    valueProvider valueprovider.Provider,
    logger logger.Logger,
) *Worker {
    return &Worker{
        exampleRepo:   exampleRepo,
        cacheRepo:     cacheRepo,
        kafkaQueue:    kafkaQueue,
        valueProvider: valueProvider,
        logger:        logger,
    }
}

func (w *Worker) Start() {
    ctx := context.Background()
    w.logger.Info(ctx, "worker started, reading tasks from kafka...")
    for {
        // читаем сообщение
        jsonData, message, err := w.kafkaQueue.ReadTask()
        if err != nil {
            w.logger.Error(ctx, "failed to unmarshal task", "error", err)
            continue
        }

        var task models.Task
        if err := json.Unmarshal(jsonData, &task); err != nil {
            w.logger.Error(ctx, fmt.Sprintf("failed to unmarshal task: %v", err))
            continue
        }

        // обрабатываем таск
        result, err := w.ProcessTask(ctx, task)

        // если ошибка в выражении (деление на ноль и т.п.) — пишем в бд и коммитим
        if err != nil && (err == calculator.ErrDivisionByZero || err == calculator.ErrCovertExample) {
            w.logger.Debug(ctx, "business error in task", "task Variable", task.Variable, "error", err)
            if errDB := w.exampleRepo.UpdateExampleWithError(ctx, task.ExampleID, err.Error()); errDB != nil {
                w.logger.Error(ctx, fmt.Sprintf("failed to save error to db: %v", errDB))
                // не коммитим — попробуем снова
                continue
            }
            // сохранили ошибку — можно коммитить
            if errCommit := w.kafkaQueue.Commit(message); errCommit != nil {
                w.logger.Error(ctx, fmt.Sprintf("failed to commit after error: %v", errCommit))
                continue
            }
            continue
        }

        // если ошибка не в выражении (redis, сеть и т.п.) — не коммитим, кавка повторит
        if err != nil {
            w.logger.Error(ctx, fmt.Sprintf("infra error, will retry: %v", err))
            continue
        }

        // всё ок — коммитим
        if err := w.kafkaQueue.Commit(message); err != nil {
            w.logger.Error(ctx, "failed to commit message", "error", err)
            continue
        }

        // если это финальный таск — сохраняем результат
        if task.IsFinal {
            if err := w.handleFinalTask(ctx, task, result); err != nil {
                w.logger.Error(ctx, "failed to save final result", "variable", task.Variable, "error", err)
                // не критично — результат уже в редисе
            }
        }
    }
}

// ProcessTask — выполняет вычисление и сохраняет в Redis
func (w *Worker) ProcessTask(ctx context.Context, task models.Task) (float64, error) {
    val1, err := w.valueProvider.Resolve(ctx, task.Num1)
    if err != nil {
        return 0, fmt.Errorf("resolve num1 (%s): %w", task.Num1, err)
    }
    val2, err := w.valueProvider.Resolve(ctx, task.Num2)
    if err != nil {
        return 0, fmt.Errorf("resolve num2 (%s): %w", task.Num2, err)
    }

    calc := calculator.NewNode(val1, val2, task.Sign)
    result, err := calc.Calculate()
    if err != nil {
        return 0, err
    }

    // Сохраняем в Redis
    if err := w.cacheRepo.SetResult(ctx, task.Variable, result); err != nil {
        return 0, fmt.Errorf("save result to Redis: %w", err)
    }

    w.logger.Info(ctx, "task processed", "num1", val1, "sign", task.Sign, "num2", val2, "result", result, "response", task.Variable)
    return result, nil
}

// handleFinalTask — обновляет Example в БД
func (w *Worker) handleFinalTask(ctx context.Context, task models.Task, result float64) error {
    if err := w.exampleRepo.UpdateExample(ctx, task.ExampleID, result); err != nil {
        return fmt.Errorf("update example in DB: %w", err)
    }

    w.logger.Info(ctx, "final result saved", "example", task.ExampleID, "result", result)
    return nil
}
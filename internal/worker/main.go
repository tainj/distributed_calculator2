// internal/worker/worker.go
package worker

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "github.com/tainj/distributed_calculator2/internal/models"
    repo "github.com/tainj/distributed_calculator2/internal/repository"
    "github.com/tainj/distributed_calculator2/internal/valueprovider"
    "github.com/tainj/distributed_calculator2/pkg/calculator"
    "github.com/tainj/distributed_calculator2/pkg/messaging/kafka"
)

type Worker struct {
    exampleRepo   repo.ExampleRepository
    cacheRepo     repo.VariableRepository
    kafkaQueue    kafka.TaskQueue
    valueProvider valueprovider.Provider
}

func NewWorker(
    exampleRepo repo.ExampleRepository,
    cacheRepo repo.VariableRepository,
    kafkaQueue kafka.TaskQueue,
    valueProvider valueprovider.Provider,
) *Worker {
    return &Worker{
        exampleRepo:   exampleRepo,
        cacheRepo:     cacheRepo,
        kafkaQueue:    kafkaQueue,
        valueProvider: valueProvider,
    }
}

func (w *Worker) Start() {
    log.Println("worker started, reading tasks from kafka...")
    for {
        ctx := context.Background()

        // читаем сообщение
        jsonData, message, err := w.kafkaQueue.ReadTask()
        if err != nil {
            log.Printf("failed to read task from kafka: %v", err)
            continue
        }

        var task models.Task
        if err := json.Unmarshal(jsonData, &task); err != nil {
            log.Printf("failed to unmarshal task: %v", err)
            continue
        }

        // обрабатываем таск
        result, err := w.ProcessTask(ctx, task)

        // если ошибка в выражении (деление на ноль и т.п.) — пишем в бд и коммитим
        if err != nil && (err == calculator.ErrDivisionByZero || err == calculator.ErrCovertExample) {
            log.Printf("business error in task %s: %v", task.Variable, err)
            if errDB := w.exampleRepo.UpdateExampleWithError(ctx, task.ExampleID, err.Error()); errDB != nil {
                log.Printf("failed to save error to db: %v", errDB)
                // не коммитим — попробуем снова
                continue
            }
            // сохранили ошибку — можно коммитить
            if errCommit := w.kafkaQueue.Commit(message); errCommit != nil {
                log.Printf("failed to commit after error: %v", errCommit)
                continue
            }
            continue
        }

        // если ошибка не в выражении (redis, сеть и т.п.) — не коммитим, кавка повторит
        if err != nil {
            log.Printf("infra error, will retry: %v", err)
            continue
        }

        // всё ок — коммитим
        if err := w.kafkaQueue.Commit(message); err != nil {
            log.Printf("failed to commit message: %v", err)
            continue
        }

        // если это финальный таск — сохраняем результат
        if task.IsFinal {
            if err := w.handleFinalTask(ctx, task, result); err != nil {
                log.Printf("failed to save final result %s: %v", task.Variable, err)
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

    log.Printf("Task processed: %f %s %f = %f → %s", val1, task.Sign, val2, result, task.Variable)
    return result, nil
}

// handleFinalTask — обновляет Example в БД
func (w *Worker) handleFinalTask(ctx context.Context, task models.Task, result float64) error {
    if err := w.exampleRepo.UpdateExample(ctx, task.ExampleID, result); err != nil {
        return fmt.Errorf("update example in DB: %w", err)
    }

    log.Printf("Final result saved: example=%s, result=%f", task.ExampleID, result)
    return nil
}

func (w * Worker) handleTaskWithError(ctx context.Context, task models.Task, err string) error {
    if err := w.exampleRepo.UpdateExampleWithError(ctx, task.ExampleID, err); err != nil {
        return fmt.Errorf("update example in DB: %w", err)
    }

    log.Printf("Final result saved with error: example=%s, error=%s", task.ExampleID, err)
    return nil
}
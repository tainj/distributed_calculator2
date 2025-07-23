// internal/worker/worker.go
package worker

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "github.com/tainj/distributed_calculator2/internal/models"
    "github.com/tainj/distributed_calculator2/internal/repository"
    "github.com/tainj/distributed_calculator2/internal/valueprovider"
    "github.com/tainj/distributed_calculator2/pkg/calculator"
    "github.com/tainj/distributed_calculator2/pkg/messaging/kafka"
)

type Worker struct {
    exampleRepo   repository.ExampleRepository
    cacheRepo     repository.CacheRepository
    kafkaQueue    kafka.TaskQueue
    valueProvider valueprovider.Provider
}

func NewWorker(
    exampleRepo repository.ExampleRepository,
    cacheRepo repository.CacheRepository,
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
    log.Println("Worker started, reading tasks from Kafka...")
    for {
        ctx := context.Background()

        // Читаем сообщение
        jsonData, message, err := w.kafkaQueue.ReadTask()
        if err != nil {
            log.Printf("Failed to read task from Kafka: %v", err)
            continue
        }

        var task models.Task
        if err := json.Unmarshal(jsonData, &task); err != nil {
            log.Printf("Failed to unmarshal task: %v", err)
            continue
        }

        // Обрабатываем таск
        result, err := w.ProcessTask(ctx, task)
		if err != nil {
            log.Printf("Failed to process task %s: %v", task.Variable, err)
            continue // ❌ не коммитим → Kafka повторит
        }

        // ✅ Коммитим только при успехе
        if err := w.kafkaQueue.Commit(message); err != nil {
            log.Printf("Failed to commit message: %v", err)
            continue
        }

        // ✅ Проверяем, финальный ли это таск
        if task.IsFinal {
            if err := w.handleFinalTask(ctx, task, result); err != nil {
                log.Printf("Failed to handle final task %s: %v", task.Variable, err)
                // ❌ Не коммитим? Нет, уже коммитили. Но можно не обновлять статус.
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
        return 0, fmt.Errorf("calculate %f %s %f: %w", val1, task.Sign, val2, err)
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
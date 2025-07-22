package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/tainj/distributed_calculator2/internal/models"
	repo "github.com/tainj/distributed_calculator2/internal/repository"
	"github.com/tainj/distributed_calculator2/internal/valueprovider"
	"github.com/tainj/distributed_calculator2/pkg/calculator"
	"github.com/tainj/distributed_calculator2/pkg/messaging/kafka"
)

type Worker struct {
    repo       repo.ResultRepository
    kafkaQueue kafka.TaskQueue
	valueProvider valueprovider.Provider

}

func NewWorker(repo repo.ResultRepository, kafkaQueue kafka.TaskQueue, valueProvider valueprovider.Provider) *Worker {
    return &Worker{repo: repo, kafkaQueue: kafkaQueue, valueProvider: valueProvider}
}

func (w *Worker) Start() {
	for {
		// создаём контекст
		ctx := context.Background()

		// читаем задачу из Kafka
		jsonData, message, err := w.kafkaQueue.ReadTask()
		if err != nil {
			log.Printf("Failed to read task from Kafka: %v\n", err)
			continue
		}

		// десериализуем JSON
		var task models.Task
		if err := json.Unmarshal(jsonData, &task); err != nil {
			log.Printf("Failed to unmarshal task: %v\n", err)
			continue
		}

		err = w.ProcessTask(ctx, task)
		if err := w.kafkaQueue.Commit(message); err != nil {
			log.Printf("Failed to commit message: %v\n", err)
			continue
		}

		// подтверждаем обработку сообщения
		if err := w.kafkaQueue.Commit(message); err != nil {
			log.Printf("Failed to commit message: %v\n", err)
			continue
		}
	}
}

func (w *Worker) ProcessTask(ctx context.Context, task models.Task) error {
	num1, err := w.valueProvider.Resolve(ctx, task.Num1)
    if err != nil {
        return err
    }
    num2, err := w.valueProvider.Resolve(ctx, task.Num2)
    if err != nil {
        return err
    }

	resolvedTask := calculator.NewNode(num1, num2, task.Sign)

    result, err := resolvedTask.Calculate()
    if err != nil {
        return err
    }

    return w.repo.SaveResult(ctx, task.Variable, result)
}
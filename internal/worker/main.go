package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	repo "github.com/tainj/distributed_calculator2/internal/repository"
	"github.com/tainj/distributed_calculator2/kafka"
	"github.com/tainj/distributed_calculator2/pkg/calculator"
)

type Worker struct {
    repo       *repo.CalculatorRepository
    kafkaQueue *kafka.KafkaQueue
}

func NewWorker(repo *repo.CalculatorRepository, kafkaQueue *kafka.KafkaQueue) *Worker {
    return &Worker{repo: repo, kafkaQueue: kafkaQueue}
}

func (w *Worker) Start() {
	for {
		// Читаем задачу из Kafka
		jsonData, message, err := w.kafkaQueue.ReadTask()
		if err != nil {
			log.Printf("Failed to read task from Kafka: %v\n", err)
			continue
		}

		// Десериализуем JSON
		var task calculator.MathExample
		if err := json.Unmarshal(jsonData, &task); err != nil {
			log.Printf("Failed to unmarshal task: %v\n", err)
			continue
		}

		// Вычисляем результат
		result, err := task.Calculate(w.repo.Cache)
		if err != nil {
			log.Printf("Failed to calculate task: %v\n", err)
			continue
		}

		// Создаём контекст
		ctx := context.Background()

		// Сохраняем результат в Redis
		if err := w.repo.Cache.SetByKey(ctx, fmt.Sprintf("user:1:variable:%s", task.Variable), result); err != nil {
			log.Printf("Failed to save result to Redis: %v\n", err)
			continue
		}

		// Подтверждаем обработку сообщения
		if err := w.kafkaQueue.Commit(message); err != nil {
			log.Printf("Failed to commit message: %v\n", err)
			continue
		}

		log.Printf("Task processed successfully. Variable: %s, Result: %s\n", task.Variable, result)
	}
}
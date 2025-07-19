// worker/worker.go
package worker

import (
  "encoding/json"
  "github.com/Anton6896/distributed_calculator2/internal/repository"
  "github.com/Anton6896/distributed_calculator2/kafka"
)

type Worker struct {
  repo       *repository.CCalculatorRepository
  kafkaQueue *kafka.KafkaQueue
}

func NewWorker(repo *repository.CCalculatorRepository, kafkaQueue *kafka.KafkaQueue) *Worker {
  return &Worker{repo: repo, kafkaQueue: kafkaQueue}
}

type Task struct {
  Expression string `json:"expression"`
}

func (w *Worker) Start() {
  for {
    // Читаем задачу из Kafka
    jsonData, err := w.kafkaQueue.ReadTask()
    if err != nil {
      continue
    }

    // Десериализуем JSON
    var task Task
    if err := json.Unmarshal(jsonData, &task); err != nil {
      continue
    }

    // Вычисли результат
    result := evaluateExpression(task.Expression)

    // Сохраняем результат в Redis
    if err := w.repo.SaveResult(task.Expression, result); err != nil {
      continue
    }
  }
}

func evaluateExpression(expression string) string {
  // Логика вычисления (пока заглушка)
  return "42"
}
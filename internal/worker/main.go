package worker

import (
    "encoding/json"
    repo "github.com/tainj/distributed_calculator2/internal/repository"
    "github.com/tainj/distributed_calculator2/kafka"
)

type Worker struct {
    repo       *repo.CalculatorRepository
    kafkaQueue *kafka.KafkaQueue
}

func NewWorker(repo *repo.CalculatorRepository, kafkaQueue *kafka.KafkaQueue) *Worker {
    return &Worker{repo: repo, kafkaQueue: kafkaQueue}
}

type Task struct {
    Num1 string `json:"num1"`
    Num2 string `json:"num2"`
    Sign string `json:"sign"`
    Variable string `json:"variable"`
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
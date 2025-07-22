package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/segmentio/kafka-go"
)

type TaskQueue interface {
	SendTask(task interface{}) error
	ReadTask() ([]byte, kafka.Message, error)
	Commit(message kafka.Message) error
}

type KafkaQueue struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

type Config struct {
	BootstrapServers  string `env:"KAFKA_BOOTSTRAP_SERVERS" env-default:"localhost:9092,localhost:9094,localhost:9096"`
	TopicCalculations string `env:"KAFKA_TOPIC_CALCULATIONS" env-default:"calculator_tasks"`
	TopicResults      string `env:"KAFKA_TOPIC_RESULTS" env-default:"calculator_results"`
}

func NewKafkaQueue(cfg Config) (*KafkaQueue, error) {
	brokers := strings.Split(cfg.BootstrapServers, ",")
	log.Printf("Initializing KafkaQueue with brokers: %v, topic: %s\n", brokers, cfg.TopicCalculations)
	if len(brokers) == 0 {
		return nil, fmt.Errorf("no Kafka brokers specified")
	}
	return &KafkaQueue{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(brokers...),
			Topic: cfg.TopicCalculations,
		},
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   cfg.TopicCalculations,
			GroupID: "calculator_group",
		}),
	}, nil
}

func (k *KafkaQueue) SendTask(task interface{}) error {
	jsonData, err := json.Marshal(task)
	if err != nil {
		log.Printf("Failed to marshal task: %v\n", err)
		return err
	}
	log.Printf("Sending task to Kafka: %s\n", string(jsonData))
	err = k.writer.WriteMessages(context.Background(), kafka.Message{
		Value: jsonData,
	})
	if err != nil {
		log.Printf("Failed to send task to Kafka: %v\n", err)
		return err
	}
	log.Println("Task sent successfully")
	return nil
}

// ReadTask возвращает данные сообщения и само сообщение для commit
func (k *KafkaQueue) ReadTask() ([]byte, kafka.Message, error) {
	log.Println("Reading task from Kafka...")

	message, err := k.reader.ReadMessage(context.Background())
	if err != nil {
		log.Printf("Failed to read message from Kafka: %v\n", err)
		return nil, kafka.Message{}, err
	}

	log.Printf("Received message: %s (offset: %d)\n", string(message.Value), message.Offset)
	return message.Value, message, nil
}

// Commit подтверждает обработку сообщения
func (k *KafkaQueue) Commit(message kafka.Message) error {
	log.Printf("Committing message (offset: %d)\n", message.Offset)

	err := k.reader.CommitMessages(context.Background(), message)
	if err != nil {
		log.Printf("Failed to commit message: %v\n", err)
		return err
	}

	log.Println("Message committed successfully")
	return nil
}

package kafka

import (
  "context"
  "encoding/json"
  "github.com/segmentio/kafka-go"
)

type KafkaQueue struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaQueue(brokers []string, topic string) *KafkaQueue {
	return &KafkaQueue{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(brokers...),
			Topic: topic,
		},
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: "calculator_group",
		}),
	}
}

func (k *KafkaQueue) SendTask(task interface{}) error {
	jsonData, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return k.writer.WriteMessages(context.Background(), kafka.Message{
		Value: jsonData,
	})
}

func (k *KafkaQueue) ReadTask() ([]byte, error) {
	message, err := k.reader.ReadMessage(context.Background())
	if err != nil {
		return nil, err
	}

	return message.Value, nil
}
package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	fmt.Println("🔥 Сервер в разработке, но уже что-то работает!")
	
	// Просто тест подключения к Kafka (заглушка)
	if os.Getenv("KAFKA_BOOTSTRAP_SERVERS") != "" {
		log.Println("Kafka: подключение... (это фейк, потом добавишь настоящий код)")
	} else {
		log.Println("Kafka: переменные окружения не заданы")
	}

	// Бесконечный цикл, чтобы контейнер не умирал
	for {
		time.Sleep(5 * time.Second)
		log.Println("Ping...")
	}
}
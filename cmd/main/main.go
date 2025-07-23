package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	repo "github.com/tainj/distributed_calculator2/internal/repository"
	service "github.com/tainj/distributed_calculator2/internal/service"
	"github.com/tainj/distributed_calculator2/internal/transport/grpc"
	"github.com/tainj/distributed_calculator2/internal/valueprovider"
	"github.com/tainj/distributed_calculator2/internal/worker"
	"github.com/tainj/distributed_calculator2/pkg/config"
	"github.com/tainj/distributed_calculator2/pkg/db/cache"
	"github.com/tainj/distributed_calculator2/pkg/db/postgres"
	"github.com/tainj/distributed_calculator2/pkg/logger"
	"github.com/tainj/distributed_calculator2/pkg/messaging/kafka"
)

const (
	serviceName = "distributed_calculator"
)

func main() {
	ctx := context.Background()
	mainLogger := logger.New(serviceName)
	ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	
	if cfg == nil {
		panic("failed to load config")
	}

	db, err := postgres.New(cfg.Postgres)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	repoExample := repo.NewPostgresResultRepository(db)

	redis := cache.New(cfg.Redis)
	fmt.Println(redis.Client.Ping(ctx))

	// Инициализация Kafka
	kafkaQueue, err := kafka.NewKafkaQueue(cfg.Kafka)
	if err != nil {
		mainLogger.Error(ctx, "Failed to init Kafka: "+err.Error())
		panic(err)
	}

	// Создаём valueProvider
	valueProvider := valueprovider.NewRedisValueProvider(redis)

	// Создаем репозиторий и сервис
	repo := repo.NewRedisResultRepository(redis)
	srv := service.NewCalculatorService(repo, kafkaQueue, repoExample)

	// Создаем и запускаем воркер
	worker := worker.NewWorker(repo, kafkaQueue, valueProvider)
	go worker.Start() // Запускаем воркер в отдельной горутине


	grpcserver, err := grpc.New(ctx, cfg.Grpc.GRPCPort, cfg.Grpc.RestPort, srv)
	if err != nil {
		mainLogger.Error(ctx, err.Error())
		return
	}

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := grpcserver.Start(ctx); err != nil {
			mainLogger.Error(ctx, err.Error())
		}
	}()

	<-graceCh

	if err := grpcserver.Stop(ctx); err != nil {
		mainLogger.Error(ctx, err.Error())
	}
	mainLogger.Info(ctx, "Server stopped")
}

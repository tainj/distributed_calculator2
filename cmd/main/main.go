package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/tainj/distributed_calculator2/internal/auth"
    repo "github.com/tainj/distributed_calculator2/internal/repository"
    service "github.com/tainj/distributed_calculator2/internal/service"
    "github.com/tainj/distributed_calculator2/internal/transport/grpc"
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
    // базовый контекст
    ctx := context.Background()

    // инициализируем логгеры
    mainLogger := logger.New(serviceName)

    // расширяем логгер для компонентов
    kafkaLogger := mainLogger.With("component", "KafkaConsumer")
    // workerLogger := mainLogger.With("component", "Worker", "worker_id", "1")
    // httpLogger := mainLogger.With("handler", "CalculateHandler")
    ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)

    // загружаем конфиг
    cfg, err := config.LoadConfig()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if cfg == nil {
        mainLogger.Error(ctx, "failed to load config")
        os.Exit(1)
    }

    // 1. база данных
    db, err := postgres.New(cfg.Postgres)
    if err != nil {
        fmt.Println(err)
        mainLogger.Error(ctx, "failed to init postgres", "error", err)
        os.Exit(1)
    }

    // 2. кэш (redis)
    redis := cache.New(cfg.Redis, mainLogger)
    fmt.Println(redis.Client.Ping(ctx)) // проверяем соединение

    // 3. фабрика репозиториев
    factory := repo.NewRepositoryFactory(db, redis, mainLogger)

    // 4. jwt сервис — нужен для auth middleware
    jwtService := auth.NewJWTService(cfg.JWT)


    // 5. kafka — очередь задач
    kafkaQueue, err := kafka.NewKafkaQueue(cfg.Kafka, kafkaLogger)
    if err != nil {
        mainLogger.Error(ctx, "failed to init kafka: "+err.Error())
        os.Exit(1)
    }

    // 6. valueprovider — для получения переменных из redis
    // valueProvider := valueprovider.NewRedisValueProvider(redis)

    // 7. репозитории
    // variableRepo := factory.CreateVariableRepository() // для сохранения результатов
    exampleRepo := factory.CreateExampleRepository()   // для сохранения выражений
	userRepo := factory.CreateUserRepository()

    // 8. сервис калькулятора
    srv := service.NewCalculatorService(userRepo, exampleRepo, jwtService, kafkaQueue, mainLogger)

    // 9. воркер — обрабатывает задачи из kafka
    // worker := worker.NewWorker(exampleRepo, variableRepo, kafkaQueue, valueProvider, workerLogger)
    // go worker.Start() // в отдельной горутине

    // 10. grpc сервер (gRPC + REST через gateway)
    grpcServer, err := grpc.New(ctx, cfg.Grpc.GRPCPort, cfg.Grpc.RestPort, srv, jwtService)
    if err != nil {
        mainLogger.Error(ctx, err.Error())
        return
    }

    // graceful shutdown
    graceCh := make(chan os.Signal, 1)
    signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

    // запускаем сервер асинхронно
    go func() {
        if err := grpcServer.Start(ctx); err != nil {
            mainLogger.Error(ctx, err.Error())
        }
    }()

    // ждём сигнала остановки
    <-graceCh

    // останавливаем
    if err := grpcServer.Stop(ctx); err != nil {
        mainLogger.Error(ctx, err.Error())
    }
    mainLogger.Info(ctx, "server stopped")
}
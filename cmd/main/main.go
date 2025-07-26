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
    // базовый контекст
    ctx := context.Background()

    // инициализируем логгер
    mainLogger := logger.New(serviceName)
    ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)

    // загружаем конфиг
    cfg, err := config.LoadConfig()
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
    if cfg == nil {
        panic("failed to load config")
    }

    // 1. база данных
    db, err := postgres.New(cfg.Postgres)
    if err != nil {
        fmt.Println(err)
        panic(err)
    }

    // 2. кэш (redis)
    redis := cache.New(cfg.Redis)
    fmt.Println(redis.Client.Ping(ctx)) // проверяем соединение

    // 3. фабрика репозиториев
    factory := repo.NewRepositoryFactory(db, redis)

    // 4. jwt сервис — нужен для auth middleware
    jwtService := auth.NewJWTService(cfg.JWT)

    // 5. kafka — очередь задач
    kafkaQueue, err := kafka.NewKafkaQueue(cfg.Kafka)
    if err != nil {
        mainLogger.Error(ctx, "failed to init kafka: "+err.Error())
        panic(err)
    }

    // 6. valueprovider — для получения переменных из redis
    valueProvider := valueprovider.NewRedisValueProvider(redis)

    // 7. репозитории
    variableRepo := factory.CreateVariableRepository() // для сохранения результатов
    exampleRepo := factory.CreateExampleRepository()   // для сохранения выражений
	userRepo := factory.CreateUserRepository()

    // 8. сервис калькулятора
    srv := service.NewCalculatorService(userRepo, kafkaQueue, exampleRepo)

    // 9. воркер — обрабатывает задачи из kafka
    worker := worker.NewWorker(exampleRepo, variableRepo, kafkaQueue, valueProvider)
    go worker.Start() // в отдельной горутине

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
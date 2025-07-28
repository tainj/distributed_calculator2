package grpc

import (
    "context"
    "fmt"
    "log"
    "log/slog"
    "net"
    "net/http"

    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "github.com/rs/cors" // ДОБАВЬ ЭТО
    "github.com/tainj/distributed_calculator2/internal/auth"
    "github.com/tainj/distributed_calculator2/internal/transport/grpc/handlers"
    "github.com/tainj/distributed_calculator2/internal/transport/grpc/middlewares"
    client "github.com/tainj/distributed_calculator2/pkg/api"
    "github.com/tainj/distributed_calculator2/pkg/logger"
    "golang.org/x/sync/errgroup"
    "google.golang.org/grpc"
)

type Server struct {
    grpcServer *grpc.Server
    restServer *http.Server
    listener   net.Listener
}

func New(ctx context.Context,
    port, restPort int,
    service handlers.Service,
    jwtService auth.JWTService) (*Server, error) {

    lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    // Настройка gRPC сервера
    opts := []grpc.ServerOption{
        grpc.ChainUnaryInterceptor(
            handlers.ContextWithLogger(logger.GetLoggerFromCtx(ctx)),
        ),
    }

    grpcServer := grpc.NewServer(opts...)
    client.RegisterCalculatorServer(grpcServer, handlers.NewCalculatorService(service))

    // Настройка REST шлюза (grpc-gateway)
    restSrv := runtime.NewServeMux()
    if err := client.RegisterCalculatorHandlerServer(context.Background(), restSrv, handlers.NewCalculatorService(service)); err != nil {
        return nil, err
    }

    // Применяем твои middleware: логирование, auth и т.д.
    finalHandler := middlewares.Apply(restSrv,
        middlewares.LoggerProvider("calculator-gateway"),
        middlewares.AuthMiddleware(jwtService),
        middlewares.Logging(),
    )

    // 🔥 ДОБАВЛЯЕМ CORS В САМОЕ КОНЦО — ОБЯЗАТЕЛЬНО ПОСЛЕ ВСЕХ MIDDLEWARE
    // Это важно: CORS должен быть ВНЕШНИМ слоем
    corsHandler := cors.New(cors.Options{
        AllowedOrigins: []string{
            "http://localhost:5173", // Vite
            "http://localhost:3000", // CRA (если будешь использовать)
        },
        AllowedMethods: []string{
            "POST",
            "GET",
            "OPTIONS", // Обязательно!
        },
        AllowedHeaders: []string{
            "*",
        },
        ExposedHeaders: []string{
            "Content-Length",
        },
        AllowCredentials: true, // если используешь куки или Authorization
        MaxAge:           3600, // кэширование preflight
    }).Handler(finalHandler)

    // Создаём HTTP-сервер с CORS
    httpServer := &http.Server{
        Addr:    fmt.Sprintf(":%d", restPort),
        Handler: corsHandler, // ← ВАЖНО: передаём обработчик с CORS
    }

    return &Server{grpcServer, httpServer, lis}, nil
}

func (s *Server) Start(ctx context.Context) error {
    eg := errgroup.Group{}

    eg.Go(func() error {
        slog.InfoContext(ctx, "starting gRPC server",
            "port", s.listener.Addr().(*net.TCPAddr).Port,
        )
        return s.grpcServer.Serve(s.listener)
    })

    eg.Go(func() error {
        slog.InfoContext(ctx, "starting REST server",
            "port", s.restServer.Addr,
        )
        return s.restServer.ListenAndServe()
    })

    return eg.Wait()
}

func (s *Server) Stop(ctx context.Context) error {
    s.grpcServer.GracefulStop()
    l := logger.GetLoggerFromCtx(ctx)
    if l != nil {
        l.Info(ctx, "gRPC server stopped")
    }
    return s.restServer.Shutdown(ctx)
}
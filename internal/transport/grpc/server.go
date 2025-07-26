package grpc

import (
    "context"
    "fmt"
    "log"
    "log/slog"
    "net"
    "net/http"

    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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

// new создаёт gRPC + REST сервер с middleware
func New(ctx context.Context, 
    port, restPort int, 
    service handlers.Service,
    jwtService *auth.JWTService,) (*Server, error) {
    // слушаем порт для gRPC
    lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    // настраиваем gRPC с interceptor'ами
    opts := []grpc.ServerOption{
        grpc.ChainUnaryInterceptor(
            handlers.ContextWithLogger(logger.GetLoggerFromCtx(ctx)), // логгер в контекст
        ),
    }

    grpcServer := grpc.NewServer(opts...)
    client.RegisterCalculatorServer(grpcServer, handlers.NewCalculatorService(service))

    // создаём REST шлюз (grpc-gateway)
    restSrv := runtime.NewServeMux()
    if err := client.RegisterCalculatorHandlerServer(context.Background(), restSrv, handlers.NewCalculatorService(service)); err != nil {
        return nil, err
    }

    // оборачиваем в middleware: логирование, auth, etc
    finalHandler := middlewares.Apply(restSrv,
        middlewares.LoggerProvider("calculator-gateway"), // логгер для HTTP
        middlewares.AuthMiddleware(jwtService),           // проверка JWT
        middlewares.Logging(),                           // логируем запросы
    )

    // создаём HTTP сервер
    httpServer := &http.Server{
        Addr:    fmt.Sprintf(":%d", restPort), // порт для REST
        Handler: finalHandler,                 // используем обработчик с middleware
    }
    
    return &Server{grpcServer, httpServer, lis}, nil
}

// start запускает gRPC и REST серверы
func (s *Server) Start(ctx context.Context) error {
    eg := errgroup.Group{}

    // запускаем gRPC
    eg.Go(func() error {
        slog.InfoContext(ctx, "starting gRPC server",
            "port", s.listener.Addr().(*net.TCPAddr).Port,
        )
        return s.grpcServer.Serve(s.listener)
    })

    // запускаем REST
    eg.Go(func() error {
        slog.InfoContext(ctx, "starting REST server",
            "port", s.restServer.Addr,
        )
        return s.restServer.ListenAndServe()
    })

    return eg.Wait()
}

// stop останавливает серверы
func (s *Server) Stop(ctx context.Context) error {
    s.grpcServer.GracefulStop()
    l := logger.GetLoggerFromCtx(ctx)
    if l != nil {
        l.Info(ctx, "gRPC server stopped")
    }
    return s.restServer.Shutdown(ctx)
}
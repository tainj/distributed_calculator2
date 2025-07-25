package grpc

import (
    "context"
    "fmt"
    "log"
    "log/slog"
    "net"
    "net/http"

    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "github.com/tainj/distributed_calculator2/internal/transport/grpc/handlers"
    "github.com/tainj/distributed_calculator2/internal/transport/grpc/auth"
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

func New(ctx context.Context, port, restPort int, service handlers.Service) (*Server, error) {
    lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    // Создаем security handler
    securityHandler := auth.NewSecurityHandler()

    // Настраиваем gRPC сервер с interceptor'ами
    opts := []grpc.ServerOption{
        grpc.ChainUnaryInterceptor(
            handlers.ContextWithLogger(logger.GetLoggerFromCtx(ctx)),
            securityHandler.AuthInterceptor, // ← Аутентификация для gRPC
        ),
    }

    grpcServer := grpc.NewServer(opts...)
    client.RegisterCalculatorServer(grpcServer, handlers.NewCalculatorService(service))

    // Создаем gRPC-gateway
    restSrv := runtime.NewServeMux()
    if err := client.RegisterCalculatorHandlerServer(context.Background(), restSrv, handlers.NewCalculatorService(service)); err != nil {
        return nil, err
    }

    // Оборачиваем gateway в HTTP middleware
    finalHandler := middlewares.Apply(restSrv,
        middlewares.LoggerProvider("calculator-gateway"),
        middlewares.AuthMiddleware(), // ← Аутентификация для HTTP
        middlewares.Logging(),
    )

    // Создаем HTTP сервер с обернутым handler'ом
    httpServer := &http.Server{
        Addr:    fmt.Sprintf(":%d", restPort),
        Handler: finalHandler, // ← ВАЖНО: используем finalHandler, а не restSrv
    }
    
    return &Server{grpcServer, httpServer, lis}, nil
}

func (s *Server) Start(ctx context.Context) error {
    eg := errgroup.Group{}

    eg.Go(func() error {
        slog.InfoContext(ctx, "Starting gRPC server",
            "port", s.listener.Addr().(*net.TCPAddr).Port,
        )
        return s.grpcServer.Serve(s.listener)
    })

    eg.Go(func() error {
        slog.InfoContext(ctx, "Starting REST server",
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
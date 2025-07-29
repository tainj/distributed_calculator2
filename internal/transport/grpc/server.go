package grpc

import (
    "context"
    "fmt"
    "net"
    "net/http"

    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "github.com/rs/cors"
    "github.com/tainj/distributed_calculator2/internal/auth"
    "github.com/tainj/distributed_calculator2/internal/transport/grpc/handlers"
    "github.com/tainj/distributed_calculator2/internal/transport/grpc/middlewares"
    client "github.com/tainj/distributed_calculator2/pkg/api"
    "github.com/tainj/distributed_calculator2/pkg/logger"
    "golang.org/x/sync/errgroup"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
)

// сервер с grpc и rest
type Server struct {
    grpcServer *grpc.Server
    restServer *http.Server
    listener   net.Listener
}

// new создаёт новый grpc + rest сервер
func New(ctx context.Context,
    port, restPort int,
    service handlers.Service,
    jwtService auth.JWTService) (*Server, error) {

    // берём логгер из контекста
    loggerFromCtx := logger.GetLoggerFromCtx(ctx)

    lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
    if err != nil {
        loggerFromCtx.Error(ctx, "failed to listen", "error", err)
    }

    // настраиваем grpc с логированием
    opts := []grpc.ServerOption{
        grpc.ChainUnaryInterceptor(
            func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
                // начало запроса
                loggerFromCtx.Info(ctx, "grpc request started",
                    "method", info.FullMethod,
                    "request_type", fmt.Sprintf("%T", req))

                // вызываем обработчик
                resp, err := handler(ctx, req)

                // результат
                if err != nil {
                    loggerFromCtx.Error(ctx, "grpc request failed",
                        "method", info.FullMethod,
                        "error", err.Error())
                } else {
                    loggerFromCtx.Info(ctx, "grpc request completed",
                        "method", info.FullMethod)
                }

                return resp, err
            },
        ),
    }

    grpcServer := grpc.NewServer(opts...)
    calculatorService := handlers.NewCalculatorService(service)
    client.RegisterCalculatorServer(grpcServer, calculatorService)

    // rest шлюз
    restMux := runtime.NewServeMux(
        runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
            md := metadata.Pairs()
            // пробрасываем нужные заголовки
            for key, values := range r.Header {
                if key == "Authorization" || key == "Content-Type" {
                    md.Set(key, values...)
                }
            }
            return md
        }),
    )

    if err := client.RegisterCalculatorHandlerServer(context.Background(), restMux, calculatorService); err != nil {
        return nil, err
    }

    // middleware для rest с логами
    loggingMiddleware := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // начало http запроса
            loggerFromCtx.Info(ctx, "http request started",
                "method", r.Method,
                "url", r.URL.Path,
                "remote_addr", r.RemoteAddr,
                "user_agent", r.UserAgent())

            // чтобы поймать статус
            wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

            // вызываем обработчик
            next.ServeHTTP(wrapped, r)

            // завершение
            loggerFromCtx.Info(ctx, "http request completed",
                "method", r.Method,
                "url", r.URL.Path,
                "status", wrapped.statusCode)
        })
    }

    // все middleware вместе
    finalHandler := middlewares.Apply(
        loggingMiddleware(restMux),
        middlewares.LoggerProvider("calculator-gateway"),
        middlewares.AuthMiddleware(jwtService),
        middlewares.Logging(),
    )

    // cors — разрешаем фронт
    corsHandler := cors.New(cors.Options{
        AllowedOrigins: []string{
            "http://localhost:5173", // vite
            "http://localhost:3000", // create-react-app
        },
        AllowedMethods: []string{
            "POST", "GET", "OPTIONS", "PUT", "DELETE",
        },
        AllowedHeaders: []string{
            "Accept", "Content-Type", "Content-Length", "Authorization", "X-Requested-With",
        },
        ExposedHeaders: []string{
            "Content-Length",
        },
        AllowCredentials: true,
        MaxAge:           3600,
    }).Handler(finalHandler)

    // создаём http сервер
    httpServer := &http.Server{
        Addr:    fmt.Sprintf(":%d", restPort),
        Handler: corsHandler,
    }

    return &Server{grpcServer, httpServer, lis}, nil
}

// вспомогательная структура — чтобы получить статус ответа
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

// start запускает grpc и rest серверы
func (s *Server) Start(ctx context.Context) error {
    l := logger.GetLoggerFromCtx(ctx)
    eg := errgroup.Group{}

    eg.Go(func() error {
        l.Info(ctx, "starting grpc server", "port", s.listener.Addr().(*net.TCPAddr).Port)
        return s.grpcServer.Serve(s.listener)
    })

    eg.Go(func() error {
        l.Info(ctx, "starting rest server", "port", s.restServer.Addr)
        return s.restServer.ListenAndServe()
    })

    return eg.Wait()
}

// stop останавливает серверы
func (s *Server) Stop(ctx context.Context) error {
    s.grpcServer.GracefulStop()
    l := logger.GetLoggerFromCtx(ctx)
    if l != nil {
        l.Info(ctx, "grpc server stopped")
    }
    return s.restServer.Shutdown(ctx)
}
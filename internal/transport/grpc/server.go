package grpc

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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

func New(ctx context.Context, port, restPort int, service Service) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			ContextWithLogger(logger.GetLoggerFromCtx(ctx)),
		),
	}

	grpcServer := grpc.NewServer(opts...)
	client.RegisterCalculatorServer(grpcServer, NewOrderService(service))

	restSrv := runtime.NewServeMux()
	if err := client.RegisterCalculatorHandlerServer(context.Background(), restSrv, NewOrderService(service)); err != nil {
		return nil, err
	}
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", restPort),
		Handler: restSrv,
	}
	return &Server{grpcServer, httpServer, lis}, nil
}

func (s *Server) Start(ctx context.Context) error {
	eg := errgroup.Group{} // Создаём группу для управления горутинами

	eg.Go(func() error {
		// логируем запуск gRPC сервера
		slog.InfoContext(ctx, "Starting gRPC server",
			"port", s.listener.Addr().(*net.TCPAddr).Port,
		)
		return s.grpcServer.Serve(s.listener) // Запускаем gRPC сервер
	})

	eg.Go(func() error {
		// логируем запуск REST сервера
		slog.InfoContext(ctx, "Starting REST server",
			"port", s.restServer.Addr,
		)
		return s.restServer.ListenAndServe() // Запускаем REST сервер
	})

	return eg.Wait() // ожидаем завершения всех горутин и возвращаем ошибку, если она возникла
}
func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	l := logger.GetLoggerFromCtx(ctx)
	if l != nil {
		l.Info(ctx, "gRPC server stopped")
	}
	return s.restServer.Shutdown(ctx)
}

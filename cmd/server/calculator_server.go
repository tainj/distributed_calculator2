package server

import (
	"context"
	"log"
	"net"
	"os"

	"fmt"

	"github.com/tantoni228/distributed_calculator2/pkg/calculator"
	pb "github.com/tantoni228/distributed_calculator2/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.DistributedCalculatorServer // сервис из сгенерированного пакета
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Calculation(
	ctx context.Context,
	in *pb.SimpleExpression,
) (*pb.ValueSimpleExpression, error) {
	log.Println("invoked Calculation: ", in)
	// вычислим площадь и вернём ответ
	ch := make(chan calculator.Solution)
	go calculator.Calc(ch, in.Sign, int(in.A), int(in.B))
	data := <-ch
	return &pb.ValueSimpleExpression{Result: int32(data.Result),
		Letter: in.Letter, Status: int32(data.Status),
	}, nil
}

func CreateCalcServer(port int) {

	host := "localhost"

	addr := fmt.Sprintf("%s:%d", host, port)
	lis, err := net.Listen("tcp", addr) // будем ждать запросы по этому адресу

	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}

	log.Printf("tcp listener started at port: %v\n", port)
	// создадим сервер grpc
	grpcServer := grpc.NewServer()
	// объект структуры, которая содержит реализацию
	geomServiceServer := NewServer()
	// серверной части GeometryService
	pb.RegisterDistributedCalculatorServer(grpcServer, geomServiceServer)
	// зарегистрируем нашу реализацию сервера
	// запустим grpc сервер
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}

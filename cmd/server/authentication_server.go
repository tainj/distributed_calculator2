package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/mattn/go-sqlite3"
	"github.com/tantoni228/distributed_calculator2/database/connection"
	pb "github.com/tantoni228/distributed_calculator2/proto"
	"google.golang.org/grpc"
)

type AuthenticationServer struct {
	pb.AuthenticationServerServer // сервис из сгенерированного пакета
}

func NewAuthenticationServer() *AuthenticationServer {
	return &AuthenticationServer{}
}

func (s *AuthenticationServer ) Authorization(
	ctx context.Context,
	in *pb.Form,
) (*pb.AuthorizationStatus, error) {
	log.Println("Authorization: ", in)
	if err := connection.InsertUser(in.Login, in.Password); err != nil {
		log.Printf("failed invoking Authorization: %v", err)
		if (err.(sqlite3.Error)).ExtendedCode == sqlite3.ErrConstraintUnique {
			return &pb.AuthorizationStatus{Status: int32(261)}, nil
		}
		return &pb.AuthorizationStatus{Status: int32(561)}, nil
	}
	return &pb.AuthorizationStatus{Status: int32(200),
	}, nil
}

func (c *AuthenticationServer) Authentication(
	ctx context.Context, in *pb.Form,
) (*pb.GiveToken, error) {
	log.Println("Authentication: ", in)
	token, err := connection.GenerateTokenUser(in.Login, in.Password)
	if err != nil {
		if err ==  sql.ErrNoRows {
			return &pb.GiveToken{Status: int32(261), Token: token}, nil
		}
		if err.Error() == fmt.Errorf("invalid password").Error() {
			return &pb.GiveToken{Status: int32(371), Token: token}, nil
		}
		return &pb.GiveToken{Status: int32(561), Token: token}, nil
	}
	return &pb.GiveToken{Status: int32(200), Token: token}, nil
}

func CreateAuthServer(p int) {
	log.Println("starting tcp listener")
	host := "localhost"
	port := fmt.Sprint(p)

	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr) // будем ждать запросы по этому адресу

	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}

	log.Println("tcp listener started at port: ", port)
	// создадим сервер grpc
	grpcServer := grpc.NewServer()
	// объект структуры, которая содержит реализацию
	geomServiceServer := NewAuthenticationServer()
	// серверной части GeometryService
	pb.RegisterAuthenticationServerServer(grpcServer, geomServiceServer)
	// зарегистрируем нашу реализацию сервера
	// запустим grpc сервер
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}

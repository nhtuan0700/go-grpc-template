package grpc

import (
	"context"
	hellov1 "github.com/nhtuan0700/go-grpc-template/internal/generated/proto/hello/v1"
	"log"
	"net"

	"github.com/bufbuild/protovalidate-go"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
)

type Server interface {
	Start(context.Context) error
}

type server struct {
	handler hellov1.GreeterServiceServer
}

func NewServer(
	handler hellov1.GreeterServiceServer,
) Server {
	return &server{
		handler: handler,
	}
}

func (s server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		return err
	}
	defer listener.Close()

	validator, err := protovalidate.New()
	if err != nil {
		return err
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(protovalidate_middleware.UnaryServerInterceptor(validator)),
	)

	hellov1.RegisterGreeterServiceServer(server, s.handler)
	log.Println("Starting grpc server")
	return server.Serve(listener)
}

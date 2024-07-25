package grpc

import (
	"context"
	"fmt"
	hellov1 "github.com/nhtuan0700/go-grpc-template/internal/generated/proto/hello/v1"
)

type Handler struct {
	hellov1.UnimplementedGreeterServiceServer
}

func NewHandler() hellov1.GreeterServiceServer {
	return Handler{}
}

func (h Handler) SayHello(ctx context.Context, req *hellov1.SayHelloRequest) (*hellov1.SayHelloResponse, error) {
	greeter := fmt.Sprintf("Hello %s", req.Name)
	return &hellov1.SayHelloResponse{
		Message: greeter,
	}, nil
}

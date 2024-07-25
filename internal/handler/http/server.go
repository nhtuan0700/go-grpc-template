package http

import (
	"context"
	hellov1 "github.com/nhtuan0700/go-grpc-template/internal/generated/proto/hello/v1"
	"github.com/nhtuan0700/go-grpc-template/internal/handler/http/middleware"
	"log"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server interface {
	Start(ctx context.Context) error
}

type server struct {
}

func NewServer() Server {
	return server{}
}

func (s server) setGRPCGatewayHandler(ctx context.Context) (http.Handler, error) {
	grpcMux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := hellov1.RegisterGreeterServiceHandlerFromEndpoint(
		ctx,
		grpcMux,
		"127.0.0.1:8081",
		opts,
	)
	if err != nil {
		return nil, err
	}

	return grpcMux, err
}

func (s server) Start(ctx context.Context) error {
	grpcGatewayHandler, err := s.setGRPCGatewayHandler(ctx)
	if err != nil {
		return err
	}

	handler := middleware.ExampleMiddleware(middleware.CorsMiddleware(grpcGatewayHandler))

	httpServer := http.Server{
		Addr:        "0.0.0.0:8080",
		Handler:     handler,
		ReadTimeout: time.Minute,
	}

	log.Println("Starting http server")
	return httpServer.ListenAndServe()
}

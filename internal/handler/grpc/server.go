package grpc

import (
	"context"
	"net"

	"github.com/nhtuan0700/go-grpc-template/internal/config"
	hellov1 "github.com/nhtuan0700/go-grpc-template/internal/generated/proto/hello/v1"
	"github.com/nhtuan0700/go-grpc-template/internal/utils"
	"go.uber.org/zap"

	"github.com/bufbuild/protovalidate-go"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
)

type Server interface {
	Start(context.Context) error
}

type server struct {
	handler    hellov1.GreeterServiceServer
	grpcConfig config.GRPC
	logger     *zap.Logger
}

func NewServer(
	handler hellov1.GreeterServiceServer,
	grpcConfig config.GRPC,
	logger *zap.Logger,
) Server {
	return &server{
		handler:    handler,
		grpcConfig: grpcConfig,
		logger:     logger,
	}
}

func (s server) Start(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, s.logger)
	listener, err := net.Listen("tcp", s.grpcConfig.Address)
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

	logger.Info("Starting grpc server: " + s.grpcConfig.Address)
	return server.Serve(listener)
}

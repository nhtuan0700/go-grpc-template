package http

import (
	"context"
	"net/http"
	"time"

	"github.com/nhtuan0700/go-grpc-template/internal/config"
	hellov1 "github.com/nhtuan0700/go-grpc-template/internal/generated/proto/hello/v1"
	// "github.com/nhtuan0700/go-grpc-template/internal/handler/http/middleware"
	// "github.com/nhtuan0700/go-grpc-template/internal/utils"
	"go.uber.org/zap"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server interface {
	Start(ctx context.Context) error
}

type server struct {
	httpConfig config.HTTP
	grpcConfig config.GRPC
	logger     *zap.Logger
}

func NewServer(
	httpConfig config.HTTP,
	grpcConfig config.GRPC,
	logger *zap.Logger,
) Server {
	return server{
		httpConfig: httpConfig,
		grpcConfig: grpcConfig,
		logger:     logger,
	}
}

func (s server) setGRPCGatewayHandler(ctx context.Context) (http.Handler, error) {
	grpcMux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := hellov1.RegisterGreeterServiceHandlerFromEndpoint(
		ctx,
		grpcMux,
		s.grpcConfig.Address,
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

	// middlewares := []utils.Middleware{
	// 	middleware.ExampleMiddleware,
	// 	middleware.CorsMiddleware,
	// 	middleware.RequestMiddlewareWith(ctx),
	// }

	// _ = utils.AddChainingMiddleware(grpcGatewayHandler, middlewares...)

	httpServer := http.Server{
		Addr:        s.httpConfig.Address,
		Handler:     grpcGatewayHandler,
		ReadTimeout: time.Minute,
	}

	s.logger.Info("Starting http server: " + s.httpConfig.Address)
	return httpServer.ListenAndServe()
}

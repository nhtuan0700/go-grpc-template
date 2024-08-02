package http

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"time"

	"github.com/nhtuan0700/go-grpc-template/internal/config"
	hellov1 "github.com/nhtuan0700/go-grpc-template/internal/generated/proto/hello/v1"
	"github.com/nhtuan0700/go-grpc-template/internal/handler/http/middleware"
	"github.com/nhtuan0700/go-grpc-template/internal/utils"

	// "github.com/nhtuan0700/go-grpc-template/internal/handler/http/middleware"
	// "github.com/nhtuan0700/go-grpc-template/internal/utils"
	"go.uber.org/zap"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := os.ReadFile("cert/server-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}

func (s server) setGRPCGatewayHandler(ctx context.Context) (http.Handler, error) {
	grpcMux := runtime.NewServeMux()

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		return nil, err
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(tlsCredentials)}

	err = hellov1.RegisterGreeterServiceHandlerFromEndpoint(
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

	middlewares := []utils.Middleware{
		middleware.ExampleMiddleware,
		middleware.CorsMiddleware,
		middleware.RequestMiddlewareWith(ctx),
	}

	_ = utils.AddChainingMiddleware(grpcGatewayHandler, middlewares...)

	// address := s.httpConfig.Address
	address := ":https"
	httpServer := http.Server{
		Addr:        address,
		Handler:     grpcGatewayHandler,
		ReadTimeout: time.Minute,
	}

	s.logger.Info("Starting http server: " + address)
	return httpServer.ListenAndServe()
}

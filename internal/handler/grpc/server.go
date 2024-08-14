package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"github.com/nhtuan0700/go-grpc-template/internal/config"
	hellov1 "github.com/nhtuan0700/go-grpc-template/internal/generated/proto/hello/v1"
	"github.com/nhtuan0700/go-grpc-template/internal/utils"
	"go.uber.org/zap"

	"github.com/bufbuild/protovalidate-go"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		return err
	}

	server := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.UnaryInterceptor(protovalidate_middleware.UnaryServerInterceptor(validator)),
	)

	hellov1.RegisterGreeterServiceServer(server, s.handler)

	logger.Info("Starting grpc server: " + s.grpcConfig.Address)
	return server.Serve(listener)
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed client's certificate
	pemClientCA, err := os.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

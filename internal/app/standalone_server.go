package app

import (
	"context"

	"github.com/nhtuan0700/go-grpc-template/internal/config"
	"github.com/nhtuan0700/go-grpc-template/internal/handler/grpc"
	"github.com/nhtuan0700/go-grpc-template/internal/handler/http"
	"github.com/nhtuan0700/go-grpc-template/internal/utils"
)

type Server struct {
	grpcServer grpc.Server
	httpServer http.Server
}

func InitializeStandaloneServer(config config.Config) (*Server, func(), error) {
	logger, loggerCleanup, err := utils.InitializeLogger(config.Log)
	if err != nil {
		return nil, nil, err
	}

	grpcHandler := grpc.NewHandler()
	grpcServer := grpc.NewServer(grpcHandler, config.GRPC, logger)

	httpServer := http.NewServer(config.HTTP, config.GRPC, logger)

	cleanup := func ()  {
		loggerCleanup()
	}

	return &Server{
		grpcServer: grpcServer,
		httpServer: httpServer,
	}, cleanup, nil
}

func (s Server) Start() error {
	go func() {
		if err := s.grpcServer.Start(context.Background()); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := s.httpServer.Start(context.Background()); err != nil {
			panic(err)
		}
	}()

	utils.BlockUntilSignal()
	return nil
}

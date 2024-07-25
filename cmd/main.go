package main

import (
	"context"
	"github.com/nhtuan0700/go-grpc-template/internal/handler/grpc"
	"github.com/nhtuan0700/go-grpc-template/internal/handler/http"
	"log"
	"os"
	"os/signal"
)


func main() {
	go func ()  {
		grpcHandler := grpc.NewHandler()

		grpcServer := grpc.NewServer(grpcHandler)
		if err := grpcServer.Start(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	go func ()  {
		httpServer := http.NewServer()
		if err := httpServer.Start(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	<-done
}

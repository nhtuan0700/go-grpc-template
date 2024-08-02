package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nhtuan0700/go-grpc-template/internal/config"
	"go.uber.org/zap"
)

func getZapLoggerLevel(level string) zap.AtomicLevel {
	switch level {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}

func InitializeLogger(cfg config.Log) (*zap.Logger, func(), error) {
	config := zap.NewProductionConfig()
	filePath, err := getLogFilePath()
	if err != nil {
		return nil, nil, err
	}
	config.OutputPaths = append(config.OutputPaths, filePath)

	config.Level = getZapLoggerLevel(cfg.Level)

	logger, err := config.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	cleanup := func() {
		logger.Sync()
	}

	return logger, cleanup, nil
}

func LoggerWithContext(ctx context.Context, logger *zap.Logger) *zap.Logger {
	// 	requestID := uu

	return logger
}

func getLogFilePath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// fileInfo, errs := os.Stat(filename)
	// if errs != nil {
	// 	println("Error getting file information:", errs)
	// 	return
	// }

	return fmt.Sprintf("%s/app.log", wd), nil
}

package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/nhtuan0700/go-grpc-template/internal/app"
	"github.com/spf13/cobra"
)

func init() {
	var appCmd = &cobra.Command{
		Use:   "start",
		Short: "Start application",
		Long:  `This command is used for starting application`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Starting application...")
			err := godotenv.Load(".env")
			if err != nil {
				return err
			}
			app, cleanup, err := app.InitializeStandaloneServer()
			if err != nil {
				return err
			}
			defer cleanup()

			return app.Start()
		},
	}

	// appCmd.PersistentFlags().StringVar(&cfg.GRPC.Address, "grpc", ":8081", "http address")
	// appCmd.PersistentFlags().StringVar(&cfg.HTTP.Address, "http", ":8080", "grpc address")
	// appCmd.PersistentFlags().StringVar(&cfg.Log.Level, "level", "info", "log level")
	rootCmd.AddCommand(appCmd)
}

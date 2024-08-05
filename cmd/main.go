package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Version: "1.0",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Panic(err)
	}
}

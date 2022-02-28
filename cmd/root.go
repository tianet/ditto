/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var Schema string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ditto",
	Short: "Ditto is a CLI tool to generate random messages",
	Long: `Ditto is a CLI tool to generate random messages.

	For a sample schema, see https://github.com/tianet/ditto/schemas/sample.json

	Usage:
		To serve message as an rest api:
			ditto http server -s schema.json -p 8080
		To send messages to kafka (requires a running Kafka cluster):
			ditto kafka producer -s schema.json -b localhost:9092 -S 3 [-c 2] [-t sample]
	`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}

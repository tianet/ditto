package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	Port int
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Ditto transforms into an http server",
	Long:  `Ditto transforms into an http server`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)

	httpCmd.PersistentFlags().StringVarP(&Schema, "schema", "s", "", "Schema for the messages")
	httpCmd.PersistentFlags().IntVarP(&Port, "port", "p", 8080, "Port where to serve the application")
}

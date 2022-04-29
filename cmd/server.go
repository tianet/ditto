/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tianet/ditto/pkg/adapters/server"
	"github.com/tianet/ditto/pkg/application/schema"
)

var (
	Port int
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server [type]",
	Short: "Rest API for schemas provided",
	Long: `Rest API for the schemas provided.

	It supports multiple schemas if the name of the schema passed is a folder.
	It will create an endpoint for each file.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("Invalid number of parameters")
		}

		if !server.IsValidServerType(args[0]) {
			return fmt.Errorf("Server type is not supported")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		schema.NewCache()

		srv, err := server.NewServer(args[0], Port)
		if err != nil {
			panic(err)
		}

		err = filepath.WalkDir(Schema, func(path string, dir os.DirEntry, err error) error {
			if filepath.Ext(path) == ".json" {
				fields, err := schema.GetFields(path)
				if err != nil {
					return err
				}
				srv.AddEndpoint(server.GetEndpoint(path), fields)
			}
			return nil
		})

		if err != nil {
			panic(err)
		}

		srv.ListenAndServe()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().IntVarP(&Port, "port", "p", 8080, "Port used to start the server.")
}

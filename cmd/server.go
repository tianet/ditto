/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	server "github.com/tianet/ditto/pkg/http"
	"github.com/tianet/ditto/pkg/schema"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Rest API for schemas provided",
	Long: `Rest API for the schemas provided.

	It supports multiple schemas if the name of the schema passed is a folder.
	It will create an endpoint for each file.
`,
	Run: func(cmd *cobra.Command, args []string) {
		schema.NewCache()

		err := filepath.WalkDir(Schema, func(path string, dir os.DirEntry, err error) error {
			if filepath.Ext(path) == ".json" {
				fields, err := schema.GetFields(path)
				if err != nil {
					return err
				}
				endpoint := server.Endpoint{Fields: fields}
				fmt.Sprintf("Adding endpoint %s\n", server.GetEndpoint(path))
				http.HandleFunc(server.GetEndpoint(path), endpoint.Handler)
			}
			return nil
		})

		if err != nil {
			panic(err)
		}

		http.ListenAndServe(fmt.Sprintf(":%d", Port), nil)

	},
}

func init() {
	httpCmd.AddCommand(serverCmd)
}

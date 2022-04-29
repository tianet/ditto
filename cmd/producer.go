package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/tianet/ditto/pkg/adapters/encoder"
	"github.com/tianet/ditto/pkg/adapters/producer"
	"github.com/tianet/ditto/pkg/application/schema"
	"github.com/tianet/ditto/pkg/application/tools"
)

var (
	Host       string
	Encoding   string
	SchemaPath string
	Schedule   int
	Count      int
	TLS        bool
)

func validateProducerParameters() {
	fileInfo, err := os.Stat(Schema)
	if err != nil {
		panic(fmt.Errorf("Schema %s is not a valid path", Schema))
	}

	if Encoding != encoder.JSON && fileInfo.IsDir() {
		panic(fmt.Errorf("Multiple schemas is only supported when using json encoding"))
	}
}

// producerCmd represents the producer command
var producerCmd = &cobra.Command{
	Use:   "producer [type]",
	Short: "Produce messages",
	Long: `Produce messages. It requires a running Kafka deployment.


	It supports multiple schemas if the name of the schema passed is a folder.
	By default it will send the messages to a topic with the same name as the schema file.
	If the topic flag is passed, it will send all schemas to that topic.
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("Invalid number of parameters")
		}

		if !producer.IsValidProducerType(args[0]) {
			return fmt.Errorf("Producer type is not supported")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		validateProducerParameters()

		var wg sync.WaitGroup
		schema.NewCache()

		prod, err := producer.NewProducer(args[0], Host, TLS)
		if err != nil {
			panic(err)
		}

		err = filepath.WalkDir(Schema, func(path string, dir os.DirEntry, err error) error {
			if filepath.Ext(path) == ".json" {
				fields, err := schema.GetFields(path)
				if err != nil {
					return err
				}

				enc, err := encoder.NewEncoder(Encoding, SchemaPath)
				if err != nil {
					return err
				}

				wg.Add(1)
				go func(wg *sync.WaitGroup, path string, prod producer.Producer, enc encoder.Encoder) error {
					defer wg.Done()

					var destination string
					if Destination == "" {
						destination = tools.GetTopic(path)
					} else {
						destination = Destination
					}

					count := 0

					for {
						message, err := schema.GenerateMessage(fields)
						if err != nil {
							return err
						}

						content, err := enc.Marshal(message)
						if err != nil {
							return err
						}

						err = prod.SendMessage(content, destination)
						if err != nil {
							return err
						}

						count += 1
						fmt.Printf("Message %d sent to %s\n", count, destination)
						if Count != -1 && count >= Count {
							break
						}

						time.Sleep(time.Duration(Schedule) * time.Second)
					}
					return nil

				}(&wg, path, prod, enc)
			}
			return nil
		})

		if err != nil {
			panic(err)
		}
		wg.Wait()

	},
}

func init() {
	rootCmd.AddCommand(producerCmd)

	producerCmd.PersistentFlags().StringVarP(&Encoding, "encoding", "e", encoder.JSON, "Encoding used for the messages")
	producerCmd.PersistentFlags().StringVarP(&Host, "broker", "b", "", "Broker where to send the message")
	producerCmd.PersistentFlags().IntVarP(&Schedule, "schedule", "S", 30, "How often should messages be sent")
	producerCmd.PersistentFlags().IntVarP(&Count, "count", "c", -1, "How many messages it will send to each topic")
	producerCmd.PersistentFlags().StringVar(&SchemaPath, "schema-path", "", "Path to the folder that contains the avro schema. The name of the file should match the name of the ditto schema.")
	producerCmd.PersistentFlags().BoolVar(&TLS, "tls", false, "Connection uses tls or not. If active, it will NOT verify the certificate.")
}

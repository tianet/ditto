package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"

	"github.com/tianet/ditto/pkg/kafka"
	"github.com/tianet/ditto/pkg/schema"
	"github.com/tianet/ditto/pkg/tools"
)

// producerCmd represents the producer command
var producerCmd = &cobra.Command{
	Use:   "producer",
	Short: "Produce kafka messages",
	Long: `Produce kafka messages. It requires a running Kafka deployment.


	It supports multiple schemas if the name of the schema passed is a folder.
	By default it will send the messages to a topic with the same name as the schema file.
	If the topic flag is passed, it will send all schemas to that topic.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		schema.NewCache()

		producer := kafka.NewProducer(Broker)

		err := filepath.WalkDir(Schema, func(path string, dir os.DirEntry, err error) error {
			if filepath.Ext(path) == ".json" {
				fields, err := schema.GetFields(path)
				if err != nil {
					return err
				}

				wg.Add(1)
				go func(wg *sync.WaitGroup, path string, producer sarama.SyncProducer) {
					defer wg.Done()

					var topic string
					if Topic == "" {
						topic = tools.GetTopic(path)
					} else {
						topic = Topic
					}

					count := 0

					for {
						message, err := schema.GenerateMessage(fields)

						output, err := json.Marshal(message)
						if err != nil {
							panic(err)
						}
						msg := kafka.PrepareMessage(string(output), topic)
						_, _, err = producer.SendMessage(msg)
						if err != nil {
							panic(err)
						}

						fmt.Printf("New message sent to topic %s\n", topic)

						count += 1
						if Count != -1 && count >= Count {
							break
						}
						time.Sleep(time.Duration(Schedule) * time.Second)
					}

				}(&wg, path, producer)
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
	kafkaCmd.AddCommand(producerCmd)
}

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	Topic    string
	Broker   string
	Schedule int
	Count    int
)

// kafkaCmd represents the say command
var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "Ditto transforms into a Kafka service",
	Long:  `Ditto transforms into a Kafka service`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(kafkaCmd)

	kafkaCmd.PersistentFlags().StringVarP(&Topic, "topic", "t", "", "Topic where to send the message")
	kafkaCmd.PersistentFlags().StringVarP(&Broker, "broker", "b", "", "Broker where to send the message")
	kafkaCmd.PersistentFlags().StringVarP(&Schema, "schema", "s", "", "Schema for the messages")
	kafkaCmd.PersistentFlags().IntVarP(&Schedule, "schedule", "S", 30, "How often should messages be sent")
	kafkaCmd.PersistentFlags().IntVarP(&Count, "count", "c", -1, "How many messages it will send to each topic")
}

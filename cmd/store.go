package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "store metrics in Postgres",
	Run: func(cmd *cobra.Command, args []string) {
		err := Store()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func Store() error {

	// Connect to Kafka.
	r, err := Consumer()
	if err != nil {
		return fmt.Errorf("kafka problem: %w", err)
	}
	defer r.Close()

	// Read from Kafka.
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}

		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	// TODO: Write to Postgres.

	return nil
}

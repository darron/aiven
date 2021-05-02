package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/darron/aiven/site"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "store metrics in Postgres",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := Load("store")
		if err != nil {
			log.Fatal(err)
		}
		err = Store(cfg)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func Store(cfg Config) error {

	// Connect to Kafka.
	r, err := Consumer(cfg)
	if err != nil {
		return fmt.Errorf("kafka problem: %w", err)
	}
	defer r.Close()

	// Connect to Postgres
	db, err := DBConnect(cfg)
	if err != nil {
		return fmt.Errorf("postgres problem: %w", err)
	}

	log.Println("Connected to Kafka and Postgres - waiting for metrics")

	// Read from Kafka.
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println(err)
			break
		}

		// Some output so we know it's doing something.
		log.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		// Convert JSON to struct.
		metric, err := site.ExtractMetrics(m.Value)
		if err != nil {
			log.Println(err)
		}

		// Insert the metric into Postgres.
		query := `INSERT INTO metrics 
			(captured_at, address, response_time, status_code, regexp, regexp_status) 
			VALUES
			($1, $2, $3, $4, $5, $6)`
		_, err = db.Exec(query, metric.CapturedAt, metric.Address, metric.ResponseTime.Milliseconds(), metric.StatusCode, metric.Regexp, metric.RegexpStatus)
		if err != nil {
			log.Printf("SQL Error: %s\n", err)
		}
	}

	return nil
}

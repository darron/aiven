package cmd

import (
	"context"
	"log"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/darron/aiven/site"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "store metrics in Postgres",
	Run: func(cmd *cobra.Command, args []string) {
		// Get base config values
		cfg, err := Load("store")
		if err != nil {
			log.Fatal(err)
		}
		// Get complete AppConfig for DI
		app, err := GetAppConfig(cfg)
		if err != nil {
			log.Fatal(err)
		}
		err = Store(app)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// Store listens to a Kafka topic and writes the site.Metrics to Postgres
func Store(app *App) error {

	// Cleanup after Kafka
	defer app.KReader.Close()

	// Cleanup after Postgres
	defer app.DB.Close()

	log.Println("Connected to Kafka and Postgres - waiting for metrics")

	// Read from Kafka.
	for {
		m, err := app.KReader.ReadMessage(context.Background())
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
			continue
		}

		// Setup SaveToPostgres closure.
		saveFunc := func() error {
			err := SaveToPostgres(app, metric)
			if err != nil {
				return err
			}
			return nil
		}

		// Let's get a retry
		bo, cancel := GetBackoff(1*time.Second, 10*time.Second)
		defer cancel()

		// Let's actually SaveToPostgres with retries.
		err = backoff.Retry(saveFunc, bo)
		if err != nil {
			log.Println("Postgres Save Failed: ", err)
		}

	}

	return nil
}

// SaveToPostgres saves the metric to Postgres.
func SaveToPostgres(app *App, metric site.Metrics) error {
	// Insert the metric into Postgres.
	query := `INSERT INTO metrics 
			(captured_at, address, response_time, status_code, regexp, regexp_status) 
			VALUES
			($1, $2, $3, $4, $5, $6)`
	_, err := app.DB.Exec(query, metric.CapturedAt, metric.Address, metric.ResponseTime.Milliseconds(), metric.StatusCode, metric.Regexp, metric.RegexpStatus)
	if err != nil {
		return err
	}
	return nil
}

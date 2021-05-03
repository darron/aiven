package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/darron/aiven/site"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
)

var (
	gatherCmd = &cobra.Command{
		Use:   "gather",
		Short: "Gather metrics and save to Kafka",
		Run: func(cmd *cobra.Command, args []string) {
			// Get base config values.
			cfg, err := Load("gather")
			if err != nil {
				log.Fatal(err)
			}
			// Get complete AppConfig for DI
			app, err := GetAppConfig(cfg)
			if err != nil {
				log.Fatal(err)
			}
			err = Gather(app)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	debug bool
	loop  bool
)

func init() {
	gatherCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Show Debug Information")
	gatherCmd.Flags().BoolVarP(&loop, "loop", "l", true, "Loop forever by default")
}

func Gather(app *App) error {

	// Read website list from disk.
	sites, err := site.GetEntries(app.C.SitesList)
	if err != nil {
		return fmt.Errorf("GetEntries %q Error: %w", app.C.SitesList, err)
	}

	// Let's defer Kafka Closer.
	defer app.KWriter.Close()

	// Contact each website, set a reasonable timeout.
	// Send data to Kafka.
	// Lather, rinse, repeat.
	for {
		for _, eachSite := range sites {

			// Grab the metrics from each site.
			log.Printf("GetMetrics for %#v with timeout: %s\n", eachSite, app.C.HTTPTimeout)
			m, err := eachSite.GetMetrics(app.C.HTTPTimeout, app.HTTP, debug)
			if err != nil {
				fmt.Printf("GetMetrics error: %s\n", err)
				continue
			}

			// Convert the m struct to JSON.
			mJSON, err := json.Marshal(m)
			if err != nil {
				fmt.Printf("JSON marshal error: %s\n", err)
				continue
			}

			// Setup KafkaWriter func closure
			writeFunc := func() error {
				err := app.KWriter.WriteMessages(context.Background(),
					kafka.Message{
						Key:   []byte(time.Now().UTC().Format(time.RFC3339Nano)),
						Value: []byte(mJSON),
					},
				)
				if err != nil {
					return err
				}
				return nil
			}

			// Let's get a retry
			bo, cancel := GetBackoff(1*time.Second, 10*time.Second)
			defer cancel()

			// Write the message to Kafka.
			err = backoff.Retry(writeFunc, bo)
			if err != nil {
				log.Println("Kafka write Failed: ", err)
			}

		}
		// If we're not looping - just break here.
		if !loop {
			break
		}
		log.Printf("Sleeping for %s\n", app.C.GatherDelay)
		time.Sleep(app.C.GatherDelay)
	}

	return nil
}

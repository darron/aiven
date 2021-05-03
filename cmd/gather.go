package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/darron/aiven/site"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
)

var (
	gatherCmd = &cobra.Command{
		Use:   "gather",
		Short: "Gather metrics and save to Kafka",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := Load("gather")
			if err != nil {
				log.Fatal(err)
			}
			err = Gather(cfg)
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

func Gather(cfg Config) error {

	// Read website list from disk.
	sites, err := site.GetEntries(cfg.SitesList)
	if err != nil {
		return fmt.Errorf("GetEntries %q Error: %w", cfg.SitesList, err)
	}

	// Connect to Kafka.
	w, err := Producer(cfg)
	if err != nil {
		return fmt.Errorf("kafka problem: %w", err)
	}
	defer w.Close()

	// Contact each website, set a reasonable timeout.
	// Send data to Kafka.
	// Lather, rinse, repeat.
	for {
		for _, eachSite := range sites {

			// Grab the metrics from each site.
			log.Printf("GetMetrics for %#v with timeout: %s\n", eachSite, cfg.HTTPTimeout)
			m, err := eachSite.GetMetrics(cfg.HTTPTimeout, &http.Client{}, debug)
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

			// Write the message to Kafka.
			err = w.WriteMessages(context.Background(),
				kafka.Message{
					Key:   []byte(time.Now().UTC().Format(time.RFC3339Nano)),
					Value: []byte(mJSON),
				},
			)
			if err != nil {
				fmt.Printf("kafka write error: %s\n", err)
			}

		}
		// If we're not looping - just break here.
		if !loop {
			break
		}
		log.Printf("Sleeping for %s\n", cfg.GatherDelay)
		time.Sleep(cfg.GatherDelay)
	}

	return nil
}

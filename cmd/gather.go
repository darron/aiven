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
	debug          bool
	httpGetTimeout = 5 * time.Second
)

func init() {
	gatherCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Show Debug Information")
}

func Gather(cfg Config) error {
	fmt.Printf("Config: %#v\n", cfg)

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
	for _, eachSite := range sites {

		// Grab the metrics from each site.
		log.Printf("GetMetrics for %#v with timeout: %s\n", eachSite, httpGetTimeout)
		m, err := eachSite.GetMetrics(httpGetTimeout, &http.Client{}, debug)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Convert the m struct to JSON.
		mJSON, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
		}

	}

	return nil
}

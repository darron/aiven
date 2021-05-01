package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/darron/aiven/site"
	"github.com/spf13/cobra"
)

var (
	gatherCmd = &cobra.Command{
		Use:   "gather",
		Short: "Gather metrics and save to Kafka",
		Run: func(cmd *cobra.Command, args []string) {
			err := Gather(list)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	list           string
	debug          bool
	httpGetTimeout = 5 * time.Second
)

func init() {
	gatherCmd.Flags().StringVarP(&list, "list", "l", "websites", "List of websites")
	gatherCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Show Debug Information")
}

func Gather(filename string) error {
	// Read website list from disk.
	sites, err := site.GetEntries(filename)
	if err != nil {
		return fmt.Errorf("getSites %q Error: %w", filename, err)
	}

	// Connect to Kafka.

	// Contact each website, set a reasonable timeout.
	// Send data to Kafka.
	// Lather, rinse, repeat.
	for _, eachSite := range sites {
		// Grab the metrics from each site.
		_, err := eachSite.GetMetrics(httpGetTimeout, debug)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

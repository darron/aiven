package cmd

import (
	"fmt"
	"log"

	"github.com/darron/aiven/site"
	"github.com/spf13/cobra"
)

var (
	gatherCmd = &cobra.Command{
		Use:   "gather",
		Short: "Gather metrics and save to Kafka",
		Run: func(cmd *cobra.Command, args []string) {
			err := Gather(List)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	List string
)

func init() {
	gatherCmd.Flags().StringVarP(&List, "list", "l", "websites", "List of websites")
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
		fmt.Printf("%#v\n", eachSite)
	}
	return nil
}

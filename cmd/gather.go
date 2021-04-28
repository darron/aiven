package cmd

import (
	"log"
	"time"

	"github.com/spf13/cobra"
)

var gatherCmd = &cobra.Command{
	Use:   "gather",
	Short: "Gather metrics and save to Kafka",
	Run: func(cmd *cobra.Command, args []string) {
		err := Gather()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func Gather() error {
	for {
		log.Println("Gathering")
		time.Sleep(5 * time.Second)
	}
	return nil
}

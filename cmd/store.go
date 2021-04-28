package cmd

import (
	"log"
	"time"

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
	for {
		log.Println("Storing")
		time.Sleep(5 * time.Second)
	}
	return nil
}

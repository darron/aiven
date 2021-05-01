package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aiven",
	Short: "aiven homework",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Try running --help")
	},
}

func Root() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(gatherCmd)
	rootCmd.AddCommand(storeCmd)
}

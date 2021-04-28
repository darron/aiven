package cmd

import (
	"fmt"

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
	rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(gatherCmd)
	rootCmd.AddCommand(storeCmd)
}

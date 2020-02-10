package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "checkbridge",
	Short: "Checkbridge automates creating GitHub checks for CI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, world")
	},
}

// Execute is the entrypoint of the CLI application
func Execute() error {
	return rootCmd.Execute()
}

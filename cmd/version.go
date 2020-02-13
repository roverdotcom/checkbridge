package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is overridden by ldflags in dist/release.sh
var Version string = "development"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print checkbridge version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

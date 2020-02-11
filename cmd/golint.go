package cmd

import (
	"github.com/roverdotcom/checkbridge/parser"
	"github.com/spf13/cobra"
)

var golintCmd = &cobra.Command{
	Use:   "golint",
	Short: "Parse golint results",
	Run:   makeCobraCommand("golint", parser.NewGolinter),
}

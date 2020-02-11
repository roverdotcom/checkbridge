package cmd

import (
	"github.com/roverdotcom/checkbridge/parser"
	"github.com/spf13/cobra"
)

var mypyCmd = &cobra.Command{
	Use:   "mypy",
	Short: "Parse mypy results",
	Run:   makeCobraCommand("mypy", parser.NewMypy),
}

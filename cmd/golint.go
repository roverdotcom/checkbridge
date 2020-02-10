package cmd

import (
	"bufio"
	"os"

	"github.com/roverdotcom/checkbridge/parser"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var golintCmd = &cobra.Command{
	Use:   "golint",
	Short: "Parse golint results",
	Run: func(cmd *cobra.Command, args []string) {
		configureLogging(cmd)

		logrus.Debug("Parsing golint results")

		scanner := bufio.NewReader(os.Stdin)
		glParser := parser.NewGolinter(scanner)
		results, err := glParser.Run()
		if err != nil {
			logrus.WithError(err).Error("Error parsing golint results")
			return
		}
		if len(results.Annotations) == 0 {
			logrus.Debug("No annotations reported from golint parsing")
			return
		}

		logrus.Infof("Got results: %+v", results)
		os.Exit(1)
	},
}

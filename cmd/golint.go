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
		gl := parser.NewGolinter(scanner)
		results, err := gl.Run()
		if err != nil {
			logrus.WithError(err).Error("Error parsing golint results")
			return
		}

		if len(results.Annotations) == 0 {
			logrus.Info("No violations reported from golint")
			return
		}

		logrus.Infof("Got results: %+v", results)
		// TODO report the annotations to GitHub checks API
		os.Exit(1)
	},
}

package cmd

import (
	"bufio"
	"io"
	"os"

	"github.com/roverdotcom/checkbridge/parser"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type cobraRunner func(cmd *cobra.Command, args []string)
type parserFunc func(io.Reader) parser.Parser

func makeCobraCommand(name string, pfunc parserFunc) cobraRunner {
	return func(cmd *cobra.Command, args []string) {
		configureLogging(cmd)

		logrus.Debugf("Parsing %s results", name)

		scanner := bufio.NewReader(os.Stdin)
		gl := pfunc(scanner)
		results, err := gl.Run()
		if err != nil {
			logrus.WithError(err).Errorf("Error parsing %s results", name)
			os.Exit(3)
		}

		if len(results.Annotations) == 0 {
			logrus.Infof("No violations reported from %s", name)
			return
		}

		logrus.Infof("Got results: %+v", results)
		// TODO report the annotations to GitHub checks API

		if exitZero, err := cmd.Flags().GetBool("exit-zero"); err != nil {
			logrus.WithError(err).Error("Unable to read exit-zero flag")
			os.Exit(1)
		} else if !exitZero {
			logrus.Info("Exiting 1 due to issues found by tool. Pass --exit-zero to disable this behavior")
			// Exit non-zero to mark the result of the pipeline as failed since the tool found issues with the code
			os.Exit(1)
		}
	}
}

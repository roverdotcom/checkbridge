// Copyright (c) 2020 Rover.com
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "checkbridge",
	Short: "Checkbridge automates creating GitHub checks for CI",
	Run: func(cmd *cobra.Command, args []string) {
		configureLogging(cmd)
		if err := cmd.Usage(); err != nil {
			logrus.WithError(err).Error("Error showing command usage")
		}
	},
}

func configureLogging(cmd *cobra.Command) {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
	if isVerbose, err := cmd.Flags().GetBool("verbose"); err != nil {
		logrus.WithError(err).Error("Unable to read verbosity")
	} else if isVerbose {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Enabled verbose logging")
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("exit-zero", "z", false, "exit zero even when tool reports issues")

	rootCmd.AddCommand(golintCmd)
	rootCmd.AddCommand(mypyCmd)
}

// Execute is the entrypoint of the application
func Execute() error {
	return rootCmd.Execute()
}

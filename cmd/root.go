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
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "checkbridge",
	Short: "Checkbridge automates creating GitHub checks for CI",
	Run: func(cmd *cobra.Command, args []string) {
		configureLogging(viper.GetViper())
		if err := cmd.Usage(); err != nil {
			logrus.WithError(err).Error("Error showing command usage")
		}
	},
}

func configureLogging(vip *viper.Viper) {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
	if vip.GetBool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Enabled verbose logging")
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	viper.SetEnvPrefix("checkbridge")

	// Application behavior flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("exit-zero", "z", false, "exit zero even when tool reports issues")
	rootCmd.PersistentFlags().BoolP("annotate-only", "o", false, "only leave annotations, never mark check as failed")
	rootCmd.PersistentFlags().BoolP("mark-in-progress", "m", false, "mark check as in progress before parsing")
	rootCmd.PersistentFlags().StringP("details-url", "d", "", "details URL to send for check")

	// Authentication configuration flags
	rootCmd.PersistentFlags().IntP("application-id", "a", 0, "GitHub application ID (numeric)")
	rootCmd.PersistentFlags().IntP("installation-id", "i", 0, "GitHub installation ID (numeric)")
	rootCmd.PersistentFlags().StringP("private-key", "p", "", "GitHub application private key path or value")
	rootCmd.PersistentFlags().StringP("github-token", "t", "", "short-lived GitHub app token for checks auth")

	rootCmd.PersistentFlags().StringP("github-repo", "r", "", "GitHub repository (e.g. 'roverdotcom/checkbridge')")
	rootCmd.PersistentFlags().StringP("commit-sha", "c", "", "commit SHA to report status checks for")

	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Additional environment variables
	// Allow $GITHUB_TOKEN by convention
	viper.BindEnv("github-token", "GITHUB_TOKEN")
	// Allow $GITHUB_REPOSITORY for GitHub actions
	viper.BindEnv("github-repo", "GITHUB_REPOSITORY")
	viper.BindEnv("details-url", "BUILDKITE_BUILD_URL")
	viper.BindEnv("commit-sha", "GITHUB_SHA")
	viper.BindEnv("commit-sha", "BUILDKITE_COMMIT")

	// Sub-command registration
	rootCmd.AddCommand(golintCmd)
	rootCmd.AddCommand(mypyCmd)
	rootCmd.AddCommand(authCheckCommand)
	rootCmd.AddCommand(regexCmd)
	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	viper.AutomaticEnv()
}

// Execute is the entrypoint of the application
func Execute() error {
	return rootCmd.Execute()
}

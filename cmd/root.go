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

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("exit-zero", "z", false, "exit zero even when tool reports issues")
	rootCmd.PersistentFlags().BoolP("annotate-only", "o", false, "only leave annotations, never mark check as failed")

	rootCmd.PersistentFlags().IntP("application-id", "a", 0, "GitHub application ID (numeric)")
	rootCmd.PersistentFlags().IntP("installation-id", "i", 0, "GitHub installation ID (numeric)")
	rootCmd.PersistentFlags().StringP("private-key", "p", "", "GitHub application private key path or value")

	rootCmd.PersistentFlags().StringP("github-repo", "r", "", "GitHub repository (e.g. 'roverdotcom/checkbridge')")
	rootCmd.PersistentFlags().StringP("commit-sha", "c", "", "commit SHA to report status checks for")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("exit_zero", rootCmd.PersistentFlags().Lookup("exit-zero"))
	viper.BindPFlag("annotate_only", rootCmd.PersistentFlags().Lookup("annotate-only"))

	viper.BindPFlag("application_id", rootCmd.PersistentFlags().Lookup("application-id"))
	viper.BindPFlag("installation_id", rootCmd.PersistentFlags().Lookup("installation-id"))
	viper.BindPFlag("private_key", rootCmd.PersistentFlags().Lookup("private-key"))

	viper.BindPFlag("github_repo", rootCmd.PersistentFlags().Lookup("github-repo"))
	viper.BindPFlag("commit_sha", rootCmd.PersistentFlags().Lookup("commit-sha"))

	rootCmd.AddCommand(golintCmd)
	rootCmd.AddCommand(mypyCmd)
	rootCmd.AddCommand(authCheckCommand)
}

func initConfig() {
	viper.AutomaticEnv()
}

// Execute is the entrypoint of the application
func Execute() error {
	return rootCmd.Execute()
}

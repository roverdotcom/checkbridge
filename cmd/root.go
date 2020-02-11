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
		cmd.Usage()
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

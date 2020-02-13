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
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var authCheckCommand = &cobra.Command{
	Use:   "check-auth",
	Short: "Verify checkbridge is configured properly for GitHub auth",
	Run: func(cmd *cobra.Command, args []string) {
		vip := viper.GetViper()
		configureLogging(vip)
		if err := runAuthCheck(vip, os.Getenv); err != nil {
			logrus.WithError(err).Error("Auth check failed")
			os.Exit(2)
		}
	},
}

func runAuthCheck(vip *viper.Viper, env func(string) string) error {
	e := environment{
		vip: vip,
		env: env,
	}
	repo, err := newRepo(vip, os.Getenv)
	if err != nil {
		return err
	}
	token, err := e.githubToken(repo)
	if err != nil {
		return err
	}
	logrus.WithField("token", token).Info("Got auth token")
	return nil
}

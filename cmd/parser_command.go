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
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/roverdotcom/checkbridge/github"
	"github.com/roverdotcom/checkbridge/parser"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type cobraRunner func(cmd *cobra.Command, args []string)
type parserFunc func(io.Reader) parser.Parser

var defaultPerms = map[string]string{
	"checks": "write",
}

func getHeadSha() (string, error) {
	// TODO allow configuring via command line options, fall back to reading from
	// the repo checkout
	return os.Getenv("GITHUB_SHA"), nil
}

func makeCobraCommand(name string, pfunc parserFunc) cobraRunner {
	return func(cmd *cobra.Command, args []string) {
		configureLogging(cmd)
		repo, err := newRepo()
		if err != nil {
			logrus.WithError(err).Error("Unable to determine repository")
			os.Exit(3)
		}

		head, err := getHeadSha()
		if err != nil {
			logrus.WithError(err).Error("Unable to read head SHA. Cannot continue.")
			os.Exit(3)
		}

		auth := github.NewAuthProvider(os.Getenv)
		token, err := auth.GetToken(defaultPerms)
		if err != nil {
			logrus.WithError(err).Error("Unable to get GitHub token")
			os.Exit(4)
		}
		logrus.WithField("token", token).Debug("Got GitHub checks token")

		api := github.NewClient(token, repo.owner, repo.repo)

		logrus.Debugf("Parsing %s results", name)
		run := github.CheckRun{
			Status:  github.CheckStatusCompleted,
			Name:    name,
			HeadSHA: head,
		}

		scanner := bufio.NewReader(os.Stdin)
		gl := pfunc(scanner)

		result, err := gl.Run()
		if err != nil {
			logrus.WithError(err).Errorf("Error parsing %s results", name)
			run.Conclusion = github.CheckConclusionFailure

			if err := api.CreateCheck(run); err != nil {
				logrus.WithError(err).Error("Unable to create GitHub check for parse failure")
			}
			logrus.Info("Created GitHub check as failure for parse error")
			os.Exit(3)
		}

		run.Output.Summary = fmt.Sprintf("%s completed", name)
		run.Output.Title = name

		if len(result.Annotations) == 0 {
			logrus.Infof("No violations reported from %s", name)
			run.Conclusion = github.CheckConclusionSuccess
			if err := api.CreateCheck(run); err != nil {
				logrus.WithError(err).Error("Unable to create GitHub check")
				os.Exit(5)
			}
			logrus.Debug("Created github check for successful run")
			return
		}

		logrus.Infof("Got %d annotations", len(result.Annotations))

		// TODO allow neutral status
		run.Conclusion = github.CheckConclusionFailure
		run.Output = result

		if err := api.CreateCheck(run); err != nil {
			logrus.WithError(err).Error("Unable to create GitHub check")
			os.Exit(5)
		}
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

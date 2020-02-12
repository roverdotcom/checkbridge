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
	"fmt"
	"io"
	"os"

	"github.com/roverdotcom/checkbridge/github"
	"github.com/roverdotcom/checkbridge/parser"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type cobraRunner func(cmd *cobra.Command, args []string)
type parserFunc func(io.Reader) parser.Parser

var defaultPerms = map[string]string{
	"checks": "write",
}

type parseRunner struct {
	environment

	name  string
	parse parser.Parser
}

func (p parseRunner) run() {
	configureLogging(p.vip)

	repo, err := newRepo(p.vip, p.env)
	if err != nil {
		logrus.WithError(err).Error("Unable to determine repository")
		os.Exit(3)
	}

	head, err := getHeadSha(p.vip, p.env)
	if err != nil {
		logrus.WithError(err).Error("Unable to read head SHA. Cannot continue.")
		os.Exit(3)
	}

	api, err := p.apiClient(repo)
	if err != nil {
		logrus.WithError(err).Error("Unable to get GitHub token")
		os.Exit(4)
	}

	run := github.CheckRun{
		Status:  github.CheckStatusCompleted,
		Name:    p.name,
		HeadSHA: head,
	}

	logrus.Debugf("Parsing %s results", p.name)

	result, err := p.parse.Run()
	if err != nil {
		logrus.WithError(err).Errorf("Error parsing %s results", p.name)
		run.Conclusion = github.CheckConclusionFailure

		if err := api.CreateCheck(run); err != nil {
			logrus.WithError(err).Error("Unable to create GitHub check for parse failure")
		}
		logrus.Info("Created GitHub check as failure for parse error")
		os.Exit(3)
	}

	run.Output.Summary = fmt.Sprintf("%s completed", p.name)
	run.Output.Title = p.name

	if code := p.reportResults(run, result, api); code != 0 {
		os.Exit(code)
	}
}

func (p parseRunner) reportResults(run github.CheckRun, result parser.Result, api github.Client) int {
	if len(result.Annotations) == 0 {
		logrus.Infof("No violations reported from %s", p.name)
		run.Conclusion = github.CheckConclusionSuccess
		if err := api.CreateCheck(run); err != nil {
			logrus.WithError(err).Error("Unable to create GitHub check")
			return 5
		}
		logrus.Debug("Created github check for successful run")
		return 0
	}

	logrus.Infof("Got %d annotations", len(result.Annotations))

	if p.vip.GetBool("annotate_only") {
		run.Conclusion = github.CheckConclusionNeutral
	} else {
		run.Conclusion = github.CheckConclusionFailure
	}
	run.Output = result

	if err := api.CreateCheck(run); err != nil {
		logrus.WithError(err).Error("Unable to create GitHub check")
		return 5
	}

	if !p.vip.GetBool("exit_zero") {
		logrus.Info("Exiting 1 due to issues found by tool. Pass --exit-zero to disable this behavior")
		// Exit non-zero to mark the result of the pipeline as failed since the tool found issues with the code
		return 1
	}
	logrus.Debug("Successfully reported checks, exiting 0 at user request")
	return 0

}

func makeCobraCommand(name string, pfunc parserFunc) cobraRunner {
	return func(cmd *cobra.Command, args []string) {
		parse := pfunc(os.Stdin)
		runner := parseRunner{
			environment: environment{
				vip: viper.GetViper(),
				env: os.Getenv,
			},
			name:  name,
			parse: parse,
		}
		runner.run()
	}
}

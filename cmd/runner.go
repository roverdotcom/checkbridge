// Copyright (c) 2020 Rover.com
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"github.com/roverdotcom/checkbridge/github"
	"github.com/roverdotcom/checkbridge/parser"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type runner struct {
	name  string
	parse parser.Parser
	vip   *viper.Viper
	env   func(string) string
}

func (r runner) run() {
	exitZero := r.vip.GetBool("exit_zero")
	fmt.Println("Exit zero", exitZero)
	configureLogging(r.vip)

	repo, err := newRepo(r.vip, r.env)
	if err != nil {
		logrus.WithError(err).Error("Unable to determine repository")
		os.Exit(3)
	}

	head, err := getHeadSha(r.vip, r.env)
	if err != nil {
		logrus.WithError(err).Error("Unable to read head SHA. Cannot continue.")
		os.Exit(3)
	}

	api, err := r.apiClient(repo)
	if err != nil {
		logrus.WithError(err).Error("Unable to get GitHub token")
		os.Exit(4)
	}

	run := github.CheckRun{
		Status:  github.CheckStatusCompleted,
		Name:    r.name,
		HeadSHA: head,
	}

	logrus.Debugf("Parsing %s results", r.name)

	result, err := r.parse.Run()
	if err != nil {
		logrus.WithError(err).Errorf("Error parsing %s results", r.name)
		run.Conclusion = github.CheckConclusionFailure

		if err := api.CreateCheck(run); err != nil {
			logrus.WithError(err).Error("Unable to create GitHub check for parse failure")
		}
		logrus.Info("Created GitHub check as failure for parse error")
		os.Exit(3)
	}

	run.Output.Summary = fmt.Sprintf("%s completed", r.name)
	run.Output.Title = r.name

	if code := r.reportResults(run, result, api); code != 0 {
		os.Exit(code)
	}
}

func (r runner) apiClient(repo repo) (github.Client, error) {
	auth := github.NewAuthProvider(r.env)
	token, err := auth.GetToken(defaultPerms)
	if err != nil {
		return nil, err
	}
	logrus.WithField("token", token).Debug("Got GitHub checks token")

	return github.NewClient(token, repo.owner, repo.name), nil
}

func (r runner) reportResults(run github.CheckRun, result parser.Result, api github.Client) int {
	if len(result.Annotations) == 0 {
		logrus.Infof("No violations reported from %s", r.name)
		run.Conclusion = github.CheckConclusionSuccess
		if err := api.CreateCheck(run); err != nil {
			logrus.WithError(err).Error("Unable to create GitHub check")
			return 5
		}
		logrus.Debug("Created github check for successful run")
		return 0
	}

	logrus.Infof("Got %d annotations", len(result.Annotations))

	// TODO allow neutral status
	run.Conclusion = github.CheckConclusionFailure
	run.Output = result

	if err := api.CreateCheck(run); err != nil {
		logrus.WithError(err).Error("Unable to create GitHub check")
		return 5
	}

	if !r.vip.GetBool("exit_zero") {
		logrus.Info("Exiting 1 due to issues found by tool. Pass --exit-zero to disable this behavior")
		// Exit non-zero to mark the result of the pipeline as failed since the tool found issues with the code
		return 1
	}
	logrus.Debug("Successfully reported checks, exiting 0 at user request")
	return 0
}

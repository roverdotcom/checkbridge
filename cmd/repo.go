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
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type repo struct {
	owner string
	name  string
}

func (r repo) Name() string {
	return r.name
}

func (r repo) Owner() string {
	return r.owner
}

var githubRepoRegex = regexp.MustCompile("git@github.com:(.+)/(.+)")

func repoFromPath(path string) (repo, error) {
	repoParts := strings.Split(path, "/")
	if len(repoParts) != 2 {
		return repo{}, fmt.Errorf("malformed GITHUB_REPOSITORY: %s", path)
	}
	return repo{
		owner: repoParts[0],
		name:  repoParts[1],
	}, nil
}

func newRepo(c config) (repo, error) {
	passedRepo := c.GetString("github-repo")
	if passedRepo != "" {
		logrus.WithField("repo", passedRepo).Debug("Using repo from configuration")
		return repoFromPath(passedRepo)
	}

	bkRepo := c.GetString("buildkite-repo")
	if bkRepo != "" {
		logrus.WithField("repo", bkRepo).Debug("Using BUILDKITE_REPO environment value")
		match := githubRepoRegex.FindStringSubmatch(bkRepo)
		if match != nil {
			r := strings.TrimRight(match[2], ".git")
			return repo{
				owner: match[1],
				name:  r,
			}, nil
		}
	}

	return repo{}, errors.New("missing repository configuration")
}

func getHeadSha(c config) (string, error) {
	passedSha := c.GetString("commit-sha")
	if passedSha != "" {
		logrus.WithField("sha", passedSha).Debug("Using configured SHA")
		return passedSha, nil
	}
	cmd := exec.Command("git", "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// MIT License
//
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
	"regexp"
	"strings"
)

type repo struct {
	owner string
	name  string
}

var githubRepoRegex = regexp.MustCompile("git@github.com:(.+)/(.+)")

func newRepo(env func(string) string) (repo, error) {
	// TODO allow configuring via command line options, fall back to reading from
	// the repo checkout
	ghRepo := env("GITHUB_REPOSITORY")
	if ghRepo != "" {
		repoParts := strings.Split(ghRepo, "/")
		if len(repoParts) != 2 {
			return repo{}, fmt.Errorf("malformed GITHUB_REPOSITORY: %s", ghRepo)
		}
		return repo{
			owner: repoParts[0],
			name:  repoParts[1],
		}, nil
	}

	bkRepo := env("BUILDKITE_REPO")
	if bkRepo != "" {
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

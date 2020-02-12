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

package github

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// CheckClient is an interface to GitHub's API
type CheckClient interface {
	CreateCheck(CheckRun) error
}

type checkClient struct {
	client
	owner string
	repo  string
}

// NewCheckClient creates a GitHub API client for creating checks
func NewCheckClient(token string, repo Repo) CheckClient {
	return checkClient{
		client: client{
			apiBase:   apiBase,
			authToken: token,
		},
		owner: repo.Owner(),
		repo:  repo.Name(),
	}
}

func (c checkClient) checkURL() string {
	return fmt.Sprintf("repos/%s/%s/check-runs", c.owner, c.repo)
}

func (c checkClient) CreateCheck(check CheckRun) error {
	if len(check.Output.Annotations) > 50 {
		logrus.Warnf("More than 50 annotations provided (%d), only sending first 50", len(check.Output.Annotations))
		check.Output.Annotations = check.Output.Annotations[:50]
	}

	headers := map[string]string{
		"Accept": "application/vnd.github.antiope-preview+json",
	}
	postResponse := map[string]interface{}{}
	resp, err := c.postJSON(c.checkURL(), check, headers, &postResponse)
	if err != nil {
		return err
	}

	logrus.WithField("status", resp.Status).WithField("body", postResponse).Debug("Got check create response")
	if resp.StatusCode != 201 {
		return fmt.Errorf("error response from GitHub %d: %s", resp.StatusCode, postResponse)
	}
	return nil
}

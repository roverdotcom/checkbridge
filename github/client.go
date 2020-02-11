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

package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

const apiBase = "https://api.github.com"

// Client is an interface to GitHub's API
type Client interface {
	CreateCheck(CheckRun) error
}

type client struct {
	token   string
	owner   string
	repo    string
	apiBase string
}

// NewClient creates a GitHub API client for creating checks
func NewClient(token string, owner string, repo string) Client {
	return client{
		apiBase: apiBase,
		token:   token,
		owner:   owner,
		repo:    repo,
	}
}

func (c client) checkURL() string {
	// TODO make URL configurable for tests
	return fmt.Sprintf("%s/repos/%s/%s/check-runs", c.apiBase, c.owner, c.repo)
}

func (c client) CreateCheck(check CheckRun) error {
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(check); err != nil {
		return err
	}
	url := c.checkURL()
	logrus.WithField("url", url).Debug("Making HTTP request to GitHub check-runs API")

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/vnd.github.antiope-preview+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logrus.WithField("status", resp.Status).WithField("body", string(body)).Debug("Got check create response")
	if resp.StatusCode != 200 {
		return fmt.Errorf("error response from GitHub %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

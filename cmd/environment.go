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
	"github.com/roverdotcom/checkbridge/github"
	"github.com/sirupsen/logrus"
)

// config is an interface over *viper.Viper
type config interface {
	GetString(string) string
	GetBool(string) bool
}

// environment hadles fetching configuration for an environment,
// made an interface for easier testing
type environment interface {
	config() config
	githubToken(repo) (string, error)
	apiClient(repo) (github.CheckClient, error)
}

func newEnvironment(c config) environment {
	return concreteEnv{
		c: c,
	}
}

type concreteEnv struct {
	c config
}

func (ce concreteEnv) config() config {
	return ce.c
}

func (ce concreteEnv) githubToken(repo repo) (string, error) {
	auth := github.NewAuthProvider(ce.c)
	return auth.GetToken(repo, defaultPerms)
}

func (ce concreteEnv) apiClient(repo repo) (github.CheckClient, error) {
	token, err := ce.githubToken(repo)
	if err != nil {
		return nil, err
	}
	logrus.WithField("token", token).Debug("Got GitHub checks token")

	return github.NewCheckClient(token, repo), nil
}

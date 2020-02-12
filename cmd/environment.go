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
	"github.com/spf13/viper"
)

type environment struct {
	vip *viper.Viper
	env func(string) string
}

func (e environment) githubToken(repo repo) (string, error) {
	auth := github.NewAuthProvider(e.vip)
	return auth.GetToken(repo, defaultPerms)
}

func (e environment) apiClient(repo repo) (github.Client, error) {
	token, err := e.githubToken(repo)
	if err != nil {
		return nil, err
	}
	logrus.WithField("token", token).Debug("Got GitHub checks token")

	return github.NewClient(token, repo), nil
}

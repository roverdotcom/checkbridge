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
	"errors"

	"github.com/sirupsen/logrus"
)

// AuthProvider handles getting a GitHub token for the given permissions
type AuthProvider interface {
	GetToken(perms map[string]string) (string, error)
}

// ConfigProvider is an interface over *viper.Viper
type ConfigProvider interface {
	GetString(string) string
}

type githubAuth struct {
	config ConfigProvider
}

// NewAuthProvider creates an auth provider from the given environment
func NewAuthProvider(c ConfigProvider) AuthProvider {
	return githubAuth{
		config: c,
	}
}

func (g githubAuth) GetToken(perms map[string]string) (string, error) {
	if token := g.config.GetString("github_token"); token != "" {
		logrus.Debug("Using explicit GitHub token, skipping JWT exchange")
		return token, nil
	}
	return "", errors.New("No token provided or configured")
}

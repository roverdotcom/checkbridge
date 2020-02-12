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

import "errors"

type authConfig struct {
	token             string
	AppID             *string
	AppInstallationID *string
	PrivateKey        *string
}

// AuthProvider handles getting a GitHub token for the given permissions
type AuthProvider interface {
	GetToken(perms map[string]string) (string, error)
}

// EnvironProvider is an interface over os.Getenv to allow easier testing
type EnvironProvider func(string) string

type githubAuth struct {
	config authConfig
}

// NewAuthProvider creates an auth provider from the given environment
func NewAuthProvider(env EnvironProvider) AuthProvider {
	// TODO accept cobra/viper configuration params
	return githubAuth{
		config: authConfig{
			token: env("GITHUB_TOKEN"),
		},
	}
}

func (g githubAuth) GetToken(perms map[string]string) (string, error) {
	// TODO handle JWT/installation exchange
	if g.config.token == "" {
		return "", errors.New("No token provided or configured")
	}
	return g.config.token, nil
}

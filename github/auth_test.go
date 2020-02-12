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

package github_test

import (
	"testing"

	"github.com/roverdotcom/checkbridge/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mapEnvProvider struct {
	data map[string]string
}

func (m mapEnvProvider) Getenv(key string) string {
	if data, ok := m.data[key]; ok {
		return data
	}
	return ""
}

func TestGithubAuth_GetTokenFromEnv(t *testing.T) {
	mytoken := "mytoken"
	env := mapEnvProvider{
		data: map[string]string{
			"GITHUB_TOKEN": mytoken,
		},
	}
	auth := github.NewAuthProvider(env.Getenv)

	token, err := auth.GetToken(map[string]string{})
	require.NoError(t, err, "error getting token")
	assert.Equal(t, mytoken, token)
}

func TestGithubAuth_GetTokenNoEnv(t *testing.T) {
	env := mapEnvProvider{}
	auth := github.NewAuthProvider(env.Getenv)

	_, err := auth.GetToken(nil)
	assert.Error(t, err)
}

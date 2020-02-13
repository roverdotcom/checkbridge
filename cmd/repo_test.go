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
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type envStub struct {
	key   string
	value string
}

func (e envStub) get(key string) string {
	if key == e.key {
		return e.value
	}
	return ""
}

func TestNewRepo_FromGithub(t *testing.T) {
	vip := viper.New()
	vip.Set("github-repo", "foo/bar")
	repo, err := newRepo(vip, envStub{}.get)
	require.NoError(t, err)
	assert.Equal(t, "foo", repo.owner)
	assert.Equal(t, "bar", repo.name)
}

func TestNewRepo_FromBuildKite(t *testing.T) {
	env := envStub{
		"BUILDKITE_REPO",
		"git@github.com:org/with-dashes.git",
	}
	repo, err := newRepo(viper.New(), env.get)
	require.NoError(t, err)
	assert.Equal(t, "org", repo.owner)
	assert.Equal(t, "with-dashes", repo.name)
}

func TestNewRepo_MalformedBK(t *testing.T) {
	env := envStub{
		"BUILDKITE_REPO",
		"ssh://github.com:org|with-dashes.git",
	}
	_, err := newRepo(viper.New(), env.get)
	assert.Error(t, err)
}

func TestNewRepo_EmptyBK(t *testing.T) {
	env := envStub{
		"BUILDKITE_REPO",
		"",
	}
	_, err := newRepo(viper.New(), env.get)
	assert.Error(t, err)
}

func TestNewRepo_FromViper(t *testing.T) {
	env := envStub{}
	vip := viper.New()
	vip.Set("github-repo", "foo/bar")

	repo, err := newRepo(vip, env.get)
	require.NoError(t, err)
	assert.Equal(t, "foo", repo.owner)
	assert.Equal(t, "bar", repo.name)
}

func TestGetHeadSha_FromViper(t *testing.T) {
	vip := viper.New()
	vip.Set("commit-sha", "my-sha")

	sha, err := getHeadSha(vip)
	require.NoError(t, err)
	assert.Equal(t, "my-sha", sha)
}

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
	"os/exec"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRepo_MalformedRepo(t *testing.T) {
	vip := viper.New()
	vip.Set("github-repo", "foo-bar.git")
	_, err := newRepo(vip)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed")
}

func TestNewRepo_FromGithub(t *testing.T) {
	vip := viper.New()
	vip.Set("github-repo", "foo/bar")
	repo, err := newRepo(vip)
	require.NoError(t, err)
	assert.Equal(t, "foo", repo.owner)
	assert.Equal(t, "bar", repo.name)
}

func TestNewRepo_FromBuildKite(t *testing.T) {
	vip := viper.New()
	vip.Set("buildkite-repo", "git@github.com:org/with-dashes.git")
	repo, err := newRepo(vip)
	require.NoError(t, err)
	assert.Equal(t, "org", repo.owner)
	assert.Equal(t, "with-dashes", repo.name)
}

func TestNewRepo_MalformedBK(t *testing.T) {
	vip := viper.New()
	vip.Set("buildkite-repo", "ssh://github.com:org|with-dashes.git")
	_, err := newRepo(vip)
	assert.Error(t, err)
}

func TestNewRepo_EmptyBK(t *testing.T) {
	vip := viper.New()
	vip.Set("buildkite-repo", "")
	_, err := newRepo(vip)
	assert.Error(t, err)
}

func TestNewRepo_FromViper(t *testing.T) {
	vip := viper.New()
	vip.Set("github-repo", "foo/bar")

	repo, err := newRepo(vip)
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

func TestGetHeadSha_FromExec(t *testing.T) {
	_, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git unavailable")
	}
	sha, err := getHeadSha(viper.New())
	require.NoError(t, err)
	assert.NotEmpty(t, sha)
}

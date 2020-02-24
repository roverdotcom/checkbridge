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

package cmd

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCommand_OK(t *testing.T) {
	rootCmd.Run(rootCmd, nil)
}

func TestExecute_OK(t *testing.T) {
	assert.NoError(t, Execute())
}

func TestConfigureLogging_Verbose(t *testing.T) {
	vip := viper.New()
	vip.Set("verbose", true)
	configureLogging(vip)
}

func TestGetInput_BadPath(t *testing.T) {
	vip := viper.New()
	vip.Set("file", "bad/path")
	_, err := getInput(vip)
	assert.Error(t, err)
}

func TestGetInput_Stdin(t *testing.T) {
	in, err := getInput(viper.New())
	require.NoError(t, err)
	assert.Equal(t, in, os.Stdin)
}

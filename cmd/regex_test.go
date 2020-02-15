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
	"testing"

	"github.com/roverdotcom/checkbridge/parser"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegexExtractor_Invalid(t *testing.T) {
	extractor := makeExtractor(viper.New())
	_, err := extractor(nil)
	assert.Error(t, err)
}

func TestRegexExtractor_OK(t *testing.T) {
	vip := viper.New()
	vip.Set("line-pos", 2)
	vip.Set("path-pos", 1)
	vip.Set("message-pos", 3)
	vip.Set("warn", true)

	extractor := makeExtractor(vip)

	annotation, err := extractor([]string{"", "example.go", "1234", "message"})
	require.NoError(t, err)
	assert := assert.New(t)

	assert.Equal("example.go", annotation.Path)
	assert.Equal(parser.LevelWarning, annotation.Level)
	assert.Equal(1234, annotation.Line)
	assert.Equal(1234, annotation.EndLine)
	assert.Equal(0, annotation.Column)
}

func TestRegexExtractor_BadColumn(t *testing.T) {
	vip := viper.New()
	vip.Set("column-pos", 1)

	extractor := makeExtractor(vip)

	_, err := extractor([]string{"", "abcd"})
	require.Error(t, err)
}

func TestRegexExtractor_Line(t *testing.T) {
	vip := viper.New()
	vip.Set("line-pos", 1)

	extractor := makeExtractor(vip)

	_, err := extractor([]string{""})
	require.Error(t, err)
}

func TestRunRegex_BadRegex(t *testing.T) {
	vip := viper.New()
	vip.Set("regex", "[")
	assert.Equal(t, 2, runRegexCommand(vip, nil))
}

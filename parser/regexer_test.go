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

package parser_test

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
	"testing"

	"github.com/roverdotcom/checkbridge/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testRegexp = regexp.MustCompile("(.*)")

func TestRegexerRun_WithError(t *testing.T) {
	called := false

	extractWithError := func(matches []string) (parser.Annotation, error) {
		called = true
		return parser.Annotation{}, errors.New("whoops")
	}

	r := parser.NewRegexer(testRegexp, extractWithError, bytes.NewBufferString("test"))

	_, err := r.Run()

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestRegexerRun_SomeAnnotationsOK(t *testing.T) {
	regex := regexp.MustCompile("^(.*):([0-9]+): (.*)")
	extractor := func(matches []string) (parser.Annotation, error) {
		line, err := strconv.Atoi(matches[2])
		require.NoError(t, err)

		return parser.Annotation{
			Path: matches[1],
			Line: line,
		}, nil
	}

	r := parser.NewRegexer(regex, extractor, bytes.NewBufferString(`
foo/bar.go:123: message
`))

	result, err := r.Run()
	require.NoError(t, err)
	assert.Equal(t, 1, len(result.Annotations))
	assert.Equal(t, "foo/bar.go", result.Annotations[0].Path)
}

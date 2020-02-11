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
	"testing"

	"github.com/roverdotcom/checkbridge/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeGolinter(input string) parser.Parser {
	return parser.NewGolinter(bytes.NewBufferString(input))
}

func TestGolinter_ValidMatches(t *testing.T) {
	assert := assert.New(t)

	linter := makeGolinter(`
cmd/root.go:35:1: exported function Execute should have comment or be unexported
not a valid line
	`)

	results, err := linter.Run()
	require.NoError(t, err, "Error running parser")
	assert.Equal(1, len(results.Annotations))
	a := results.Annotations[0]
	assert.Equal("cmd/root.go", a.Path)
	assert.Equal(35, a.Line)
	assert.Equal(parser.LevelWarning, a.Level)
	assert.Equal("exported function Execute should have comment or be unexported", a.Message)
}

func TestGolinter_InvalidColumns(t *testing.T) {
	linter := makeGolinter(`foo/bar.go:abcd:1234: foo bar`)
	results, err := linter.Run()
	require.NoError(t, err, "Error running linter")
	assert.Equal(t, 0, len(results.Annotations))
}

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

func makeMypyLinter(input string) parser.Parser {
	return parser.NewMypy(bytes.NewBufferString(input))
}

func TestMypy_ValidMatches(t *testing.T) {
	assert := assert.New(t)

	linter := makeMypyLinter(`
main.py:6: error: Argument 1 to "main" has incompatible type "int"; expected "str"
Found 1 error in 1 file (checked 3 source files)`)
	results, err := linter.Run()
	require.NoError(t, err, "Error running parser")
	assert.Equal(1, len(results.Annotations))
	a := results.Annotations[0]
	assert.Equal("main.py", a.Path)
	assert.Equal(6, a.Line)
	assert.Equal(0, a.Column)
	assert.Equal(`Argument 1 to "main" has incompatible type "int"; expected "str"`, a.Message)
}

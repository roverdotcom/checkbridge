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

package parser

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
)

var mypyRegex = regexp.MustCompile(`^(.*):([0-9]+): (\w*): (.*)$`)

// NewMypy instantiates a mypy linter from a reader
func NewMypy(reader io.Reader) Parser {
	return NewRegexer(mypyRegex, extractMypy, reader)
}

func extractMypy(match []string) (Annotation, error) {
	line, err := strconv.Atoi(match[2])
	if err != nil {
		return Annotation{}, fmt.Errorf("parse line %s: %w", match[2], err)
	}

	return Annotation{
		Path:    match[1],
		Level:   LevelError,
		Line:    line,
		EndLine: line,
		Message: match[4],
	}, nil
}

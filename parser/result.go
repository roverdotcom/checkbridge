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

// Level represents an annotation level
type Level string

const (
	// LevelWarning is the warning level
	LevelWarning Level = "warning"
	// LevelError is the error level
	LevelError Level = "failure"
)

// Annotation represents a line-level annotation
type Annotation struct {
	Path    string `json:"path"`
	Line    int    `json:"start_line"`
	Column  int    `json:"column"`
	Message string `json:"message"`
	Level   Level  `json:"annotation_level"`
}

// Result holds the output of a parser
type Result struct {
	Annotations []Annotation `json:"annotations"`
	Title       string       `json:"title"`
	Summary     string       `json:"summary"`
}

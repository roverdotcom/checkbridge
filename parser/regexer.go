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
	"bufio"
	"io"
	"regexp"

	"github.com/sirupsen/logrus"
)

// AnnotationExtracter converts regex matches to annotations
type AnnotationExtracter func(matches []string) (Annotation, error)

type regexer struct {
	regex     *regexp.Regexp
	extracter AnnotationExtracter
	reader    io.Reader
}

// NewRegexer creates a Parser from a regex and an extraction func
func NewRegexer(regex *regexp.Regexp, extracter AnnotationExtracter, reader io.Reader) Parser {
	return regexer{
		reader:    reader,
		regex:     regex,
		extracter: extracter,
	}
}

func (r regexer) Run() (Result, error) {
	scanner := bufio.NewScanner(r.reader)
	annotations := []Annotation{}
	for scanner.Scan() {
		line := scanner.Text()
		match := r.regex.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		if a, err := r.extracter(match); err != nil {
			logrus.WithError(err).Errorf("Unable to extract annotation from line: %s", line)
		} else {
			annotations = append(annotations, a)
		}
	}

	if err := scanner.Err(); err != nil {
		logrus.WithError(err).Error("Error reading stdin")
		return Result{}, err
	}

	return Result{
		Annotations: annotations,
	}, nil
}

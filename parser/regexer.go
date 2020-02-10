package parser

import (
	"bufio"
	"io"
	"regexp"

	"github.com/sirupsen/logrus"
)

type annotationExtracter func(matches []string) ([]Annotation, error)

type regexer struct {
	regex     *regexp.Regexp
	extracter annotationExtracter
	reader    io.Reader
}

// NewRegexer creates a Parser from a regex and an extraction func
func NewRegexer(regex *regexp.Regexp, extracter annotationExtracter, reader io.Reader) Parser {
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
			annotations = append(annotations, a...)
		}
	}

	if err := scanner.Err(); err != nil {
		logrus.WithError(err).Error("Error reading stdin")
		return Result{}, err
	}

	// TODO overall message
	return Result{
		Annotations: annotations,
	}, nil
}

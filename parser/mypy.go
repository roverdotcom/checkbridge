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
		Message: match[4],
	}, nil
}

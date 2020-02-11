package parser

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
)

var golintRegex = regexp.MustCompile(`^(.*):([0-9]+):([0-9]+): (.*)$`)

// NewGolinter instantiates a Golinter from a reader
func NewGolinter(reader io.Reader) Parser {
	return NewRegexer(golintRegex, extractGolint, reader)
}

func extractGolint(match []string) (Annotation, error) {
	line, err := strconv.Atoi(match[2])
	if err != nil {
		return Annotation{}, fmt.Errorf("parse line %s: %w", match[2], err)
	}
	column, err := strconv.Atoi(match[3])
	if err != nil {
		return Annotation{}, fmt.Errorf("parse column %s: %w", match[3], err)
	}

	return Annotation{
		Path:    match[1],
		Level:   LevelWarning,
		Line:    line,
		Column:  column,
		Message: match[4],
	}, nil
}

package parser

import (
	"bufio"
	"io"
	"regexp"
	"strconv"

	"github.com/sirupsen/logrus"
)

var golintRegex = regexp.MustCompile(`^(.*):([0-9]+):([0-9]+): (.*)$`)

// Golinter parses the output of `golint`
type Golinter struct {
	reader io.Reader
}

// NewGolinter instantiates a Golinter from a reader
func NewGolinter(reader io.Reader) Golinter {
	return Golinter{
		reader: reader,
	}
}

// Run scans the input and parses into Results
func (g Golinter) Run() (Result, error) {
	scanner := bufio.NewScanner(g.reader)
	annotations := []Annotation{}
	for scanner.Scan() {
		match := golintRegex.FindStringSubmatch(scanner.Text())
		if match == nil {
			continue
		}
		line, err := strconv.Atoi(match[2])
		if err != nil {
			logrus.WithError(err).Errorf("Unable to parse line number: %s", match[2])
			continue
		}
		column, err := strconv.Atoi(match[3])
		if err != nil {
			logrus.WithError(err).Errorf("Unable to parse column number: %s", match[3])
			continue
		}
		annotations = append(annotations, Annotation{
			Path:   match[1],
			Level:  LevelWarning,
			Line:   line,
			Column: column,
		})
	}

	if err := scanner.Err(); err != nil {
		logrus.WithError(err).Error("Error reading stdin")
		return Result{}, err
	}

	// TODO
	return Result{
		Annotations: annotations,
	}, nil
}

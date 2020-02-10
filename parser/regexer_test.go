package parser_test

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"testing"

	"github.com/roverdotcom/checkbridge/parser"
	"github.com/stretchr/testify/assert"
)

var testRegexp = regexp.MustCompile("(.*)")

func TestRegexerRun_WithError(t *testing.T) {
	called := false

	extractWithError := func(matches []string) ([]parser.Annotation, error) {
		called = true
		return nil, errors.New("whoops")
	}

	r := parser.NewRegexer(testRegexp, extractWithError, bufio.NewReader(bytes.NewBufferString("test")))

	_, err := r.Run()

	assert.NoError(t, err)
	assert.True(t, called)
}

package parser_test

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
	"testing"

	"github.com/roverdotcom/checkbridge/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testRegexp = regexp.MustCompile("(.*)")

func TestRegexerRun_WithError(t *testing.T) {
	called := false

	extractWithError := func(matches []string) (parser.Annotation, error) {
		called = true
		return parser.Annotation{}, errors.New("whoops")
	}

	r := parser.NewRegexer(testRegexp, extractWithError, bytes.NewBufferString("test"))

	_, err := r.Run()

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestRegexerRun_SomeAnnotationsOK(t *testing.T) {
	regex := regexp.MustCompile("^(.*):([0-9]+): (.*)")
	extractor := func(matches []string) (parser.Annotation, error) {
		line, err := strconv.Atoi(matches[2])
		require.NoError(t, err)

		return parser.Annotation{
			Path: matches[1],
			Line: line,
		}, nil
	}

	r := parser.NewRegexer(regex, extractor, bytes.NewBufferString(`
foo/bar.go:123: message
`))

	result, err := r.Run()
	require.NoError(t, err)
	assert.Equal(t, 1, len(result.Annotations))
	assert.Equal(t, "foo/bar.go", result.Annotations[0].Path)
}

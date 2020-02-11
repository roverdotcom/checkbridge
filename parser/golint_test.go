package parser_test

import (
	"bytes"
	"testing"

	"github.com/roverdotcom/checkbridge/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeGolinter(input string) parser.Parser {
	return parser.NewGolinter(bytes.NewBufferString(input))
}

func TestGolinter_ValidMatches(t *testing.T) {
	assert := assert.New(t)

	linter := makeGolinter(`
cmd/root.go:35:1: exported function Execute should have comment or be unexported
not a valid line
	`)

	results, err := linter.Run()
	require.NoError(t, err, "Error running parser")
	assert.Equal(1, len(results.Annotations))
	a := results.Annotations[0]
	assert.Equal("cmd/root.go", a.Path)
	assert.Equal(35, a.Line)
	assert.Equal(parser.LevelWarning, a.Level)
	assert.Equal("exported function Execute should have comment or be unexported", a.Message)
}

func TestGolinter_InvalidColumns(t *testing.T) {
	linter := makeGolinter(`foo/bar.go:abcd:1234: foo bar`)
	results, err := linter.Run()
	require.NoError(t, err, "Error running linter")
	assert.Equal(t, 0, len(results.Annotations))
}

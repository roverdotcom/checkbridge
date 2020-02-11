package parser_test

import (
	"bytes"
	"testing"

	"github.com/roverdotcom/checkbridge/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeMypyLinter(input string) parser.Parser {
	return parser.NewMypy(bytes.NewBufferString(input))
}

func TestMypy_ValidMatches(t *testing.T) {
	assert := assert.New(t)

	linter := makeMypyLinter(`
main.py:6: error: Argument 1 to "main" has incompatible type "int"; expected "str"
Found 1 error in 1 file (checked 3 source files)`)
	results, err := linter.Run()
	require.NoError(t, err, "Error running parser")
	assert.Equal(1, len(results.Annotations))
	a := results.Annotations[0]
	assert.Equal("main.py", a.Path)
	assert.Equal(6, a.Line)
	assert.Equal(0, a.Column)
	assert.Equal(`Argument 1 to "main" has incompatible type "int"; expected "str"`, a.Message)
}

package parser

// Level represents an annotation level
type Level string

const (
	// LevelWarning is the warning level
	LevelWarning Level = "warning"
	// LevelError is the error level
	LevelError Level = "error"
)

// Annotation represents a line-level annotation
type Annotation struct {
	Path    string
	Line    int
	Column  int
	Message string
	Level   Level
}

// Result holds the output of a parser
type Result struct {
	Annotations []Annotation
}

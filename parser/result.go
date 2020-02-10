package parser

// Level represents an annotation level
type Level string

const (
	// LevelWarning is the warning level
	LevelWarning Level = "warning"
	// LevelError is the error level
	LevelError Level = "failure"
)

// Annotation represents a line-level annotation
type Annotation struct {
	Path    string `json:"path"`
	Line    int    `json:"start_line"`
	Column  int    `json:"column"`
	Message string `json:"message"`
	Level   Level  `json:"annotation_level"`
}

// Result holds the output of a parser
type Result struct {
	Annotations []Annotation `json:"annotations"`
	Title       string       `json:"title"`
	Summary     string       `json:"summary"`
}

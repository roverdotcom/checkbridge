package parser

// Parser is the general interface to a result parser
type Parser interface {
	Run() (Result, error)
}

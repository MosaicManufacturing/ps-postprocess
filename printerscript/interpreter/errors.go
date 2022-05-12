package interpreter

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type psError struct {
	message string
	line    int
	col     int
}

type SyntaxError struct {
	psError
}

func NewSyntaxError(message string, line, col int) *SyntaxError {
	return &SyntaxError{psError{message, line, col}}
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("syntax error at line %d:%d: %s", e.line, e.col, e.message)
}

type SyntaxErrorListener struct {
	*antlr.DefaultErrorListener
	error *SyntaxError
}

func NewSyntaxErrorListener() *SyntaxErrorListener {
	return new(SyntaxErrorListener)
}

func (l *SyntaxErrorListener) SyntaxError(_ antlr.Recognizer, _ interface{}, line, column int, msg string, _ antlr.RecognitionException) {
	if l.error == nil {
		l.error = NewSyntaxError(msg, line, column)
	}
}

func (l *SyntaxErrorListener) GetError() *SyntaxError {
	return l.error
}

type RuntimeError struct {
	psError
}

func NewRuntimeError(message string, line, col int) *RuntimeError {
	return &RuntimeError{psError{message, line, col}}
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("runtime error at line %d:%d: %s", e.line, e.col, e.message)
}

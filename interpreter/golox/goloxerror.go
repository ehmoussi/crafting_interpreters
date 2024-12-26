package golox

import "fmt"

type SyntaxError struct {
	line    int
	message string
}

func NewSyntaxError(line int, message string) *SyntaxError {
	return &SyntaxError{line: line, message: message}
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("[line %d] SYNTAX ERROR: %s\n", e.line, e.message)
}

type SyntaxErrors struct {
	errors []*SyntaxError
}

func NewSyntaxErrors(errors ...*SyntaxError) *SyntaxErrors {
	return &SyntaxErrors{errors: errors}
}

func (e *SyntaxErrors) Error() string {
	message := ""
	for _, err := range e.errors {
		message += err.Error()
	}
	return message
}

type RuntimeError struct {
	token   *Token
	message string
}

func NewRuntimeError(token *Token, message string) *RuntimeError {
	return &RuntimeError{token, message}
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] RUNTIME ERROR: %s\n", e.token.line, e.message)
}

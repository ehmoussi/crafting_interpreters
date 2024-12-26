package golox

import "fmt"

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   any
	line      int
}

func NewToken(tokenType TokenType, lexeme string, literal any, line int) *Token {
	return &Token{
		tokenType: tokenType, lexeme: lexeme, literal: literal, line: line,
	}
}

func (t *Token) ToString() string {
	if t.literal != nil {
		return fmt.Sprintf("%s %q %s", t.tokenType, t.lexeme, t.literal)
	} else {
		return fmt.Sprintf("%s %q", t.tokenType, t.lexeme)
	}
}

package golox

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

type Scanner struct {
	source string
	tokens []*Token

	start   int
	current int
	line    int
}

func NewScanner(source string, tokenCapacity int) *Scanner {
	return &Scanner{
		source: source, tokens: make([]*Token, 0, tokenCapacity), start: 0, current: 0, line: 1,
	}
}

func (s *Scanner) scanTokens() ([]*Token, error) {
	errs := make([]*SyntaxError, 0, 10)
	for {
		if s.isAtEnd() {
			break
		}
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			var syntaxErr *SyntaxError
			if errors.As(err, &syntaxErr) {
				errs = append(errs, syntaxErr)
			}
		}
	}
	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
	if len(errs) > 0 {
		return s.tokens, NewSyntaxErrors(errs...)
	}
	return s.tokens, nil
}

// Scan the current character and eventually the next characters to determine
// its type add it to the list of tokens
func (s *Scanner) scanToken() error {
	c := s.next()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		if s.nextMatch('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.nextMatch('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '>':
		if s.nextMatch('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '<':
		if s.nextMatch('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	// Comments
	case '/':
		if s.nextMatch('/') {
			for {
				if s.isAtEndLine() {
					// "//" comment only one line
					break
				} else {
					// Ignore all the characters in the line after //
					s.next()
				}
			}
		} else {
			s.addToken(SLASH)
		}
	// Ignored characters
	case ' ', '\r', '\t':
		break
	// New line
	case '\n':
		s.line += 1
	// String
	case '"':
		if s.nextString() {
			// Ignore the surrounding quotes
			s.addTokenWithLiteral(STRING, s.source[s.start+1:s.current-1])
		} else {
			return NewSyntaxError(s.line, "Unterminated string")
		}
	default:
		if s.isDigit(c) {
			s.nextNumber()
			ns := s.source[s.start:s.current]
			n, err := strconv.ParseFloat(ns, 64)
			if err != nil {
				log.Fatalf("Failed to parse the float %s", ns)
			}
			s.addTokenWithLiteral(NUMBER, n)
		} else if s.isAlpha(c) {
			s.nextIdentifier()
			text := s.source[s.start:s.current]
			keywordToken, ok := Keywords[text]
			if ok {
				s.addToken(keywordToken)
			} else {
				s.addToken(IDENTIFIER)
			}
		} else {
			return NewSyntaxError(s.line, fmt.Sprintf("Unexpected character: %q", c))
		}
	}
	return nil
}

func (s *Scanner) nextIdentifier() {
	for {
		if s.isAlphaNumeric(s.peek()) {
			s.next()
		} else {
			break
		}
	}
}

func (s *Scanner) nextNumber() bool {
	isConsumingFractional := false
	for {
		c := s.peek()
		if s.isDigit(c) {
			s.next()
		} else if c == '.' && s.isDigit(s.peekNext()) && !isConsumingFractional {
			s.next()
			isConsumingFractional = true
		} else {
			break
		}
	}
	return true
}

func (s *Scanner) nextString() bool {
	for {
		c := s.peek()
		if c == '"' || s.isAtEnd() {
			if c == '"' {
				s.next()
			} else {
				return false
			}
			break
		}
		if c == '\n' {
			s.line += 1
		}
		s.next()
	}
	return true
}

// Consume the next token only if it matches the expected character
func (s *Scanner) nextMatch(expected byte) bool {
	if s.isAtEnd() || (s.source[s.current] != expected) {
		return false
	} else {
		s.current += 1
		return true
	}
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, text, literal, s.line))
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

// Consume the next character of the source
func (s *Scanner) next() byte {
	c := s.source[s.current]
	s.current += 1
	return c
}

// Get the next character without consume it
func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

// Get the second next character without consume it
func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// Check if the current character has reach the end of the source
func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// Check if the current character has reach the end of the line or the end of the source
func (s *Scanner) isAtEndLine() bool {
	length := len(s.source)
	if s.current >= length || s.current+1 >= length || s.current+2 >= length {
		return true
	}
	return s.peek() == '\r' && s.peekNext() == '\n' || s.peek() == '\n'
}

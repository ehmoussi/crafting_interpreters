package golox

import (
	"errors"
)

type Parser[T any] struct {
	tokens  []*Token
	current int
}

func NewParser[T any](tokensCapacity int) *Parser[T] {
	tokens := make([]*Token, 0, tokensCapacity)
	return &Parser[T]{tokens: tokens, current: 0}
}

func (p *Parser[T]) Parse() (Expr[T], error) {
	return p.expression()
}

func (p *Parser[T]) synchronize() {
	p.next()
	for {
		if p.isAtEnd() {
			break
		}
		if p.previous().tokenType == SEMICOLON {
			return
		}
		switch p.peek().tokenType {
		case CLASS:
		case FUN:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}
		p.next()
	}
}

func (p *Parser[T]) expression() (Expr[T], error) {
	return p.equality()
}

func (p *Parser[T]) equality() (Expr[T], error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for {
		if !p.match(BANG_EQUAL, EQUAL_EQUAL) {
			break
		}
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser[T]) comparison() (Expr[T], error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}
	for {
		if !p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
			break
		}
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser[T]) term() (Expr[T], error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for {
		if !p.match(MINUS, PLUS) {
			break
		}
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser[T]) factor() (Expr[T], error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}
	for {
		if !p.match(SLASH, STAR) {
			break
		}
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser[T]) unary() (Expr[T], error) {
	if p.match(BANG, MINUS) {
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return NewUnary(p.previous(), right), nil
	} else {
		return p.primary()
	}
}

func (p *Parser[T]) primary() (Expr[T], error) {
	if p.match(TRUE) {
		return NewLiteral[T](true), nil
	} else if p.match(FALSE) {
		return NewLiteral[T](false), nil
	} else if p.match(NIL) {
		return NewLiteral[T](nil), nil
	} else if p.match(NUMBER, STRING) {
		return NewLiteral[T](p.previous().literal), nil
	} else if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RIGHT_PAREN, "expect ')' after expression")
		if err != nil {
			return nil, err
		}
		return NewGrouping(expr), nil
	}
	return nil, errors.New("expect expression")
}

func (p *Parser[T]) consume(tokenType TokenType, expectMessage string) (*Token, error) {
	if p.check(tokenType) {
		return p.next(), nil
	}
	ReportError(p.peek().line, expectMessage)
	return p.peek(), errors.New(expectMessage)
}

func (p *Parser[T]) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.next()
			return true
		}
	}
	return false
}

func (p *Parser[T]) check(tokenType TokenType) bool {
	return !p.isAtEnd() && p.tokens[p.current].tokenType == tokenType
}

func (p *Parser[T]) next() *Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *Parser[T]) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *Parser[T]) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser[T]) previous() *Token {
	return p.tokens[p.current-1]
}

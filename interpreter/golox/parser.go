package golox

import (
	"fmt"
)

type Parser[T any] struct {
	tokens  []*Token
	current int
}

func NewParser[T any](tokensCapacity int) *Parser[T] {
	tokens := make([]*Token, 0, tokensCapacity)
	return &Parser[T]{tokens: tokens, current: 0}
}

func (p *Parser[T]) Parse() ([]Stmt[T], error) {
	statements := make([]Stmt[T], 0, 100)
	for {
		if p.isAtEnd() {
			break
		}
		stmt, err := p.declaration()
		if err != nil {
			p.synchronize()
			return nil, err
		}
		statements = append(statements, stmt)
	}
	return statements, nil
}

func (p *Parser[T]) declaration() (Stmt[T], error) {
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser[T]) statement() (Stmt[T], error) {
	if p.match(PRINT) {
		return p.printStatement()
	} else if p.match(LEFT_BRACE) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return NewBlock(statements), nil
	}
	return p.expressionStatement()
}

func (p *Parser[T]) block() ([]Stmt[T], error) {
	statements := make([]Stmt[T], 0, 15)
	for {
		if p.isAtEnd() || p.check(RIGHT_BRACE) {
			break
		}
		statement, err := p.declaration()
		if err != nil {
			return statements, err
		}
		statements = append(statements, statement)
	}
	_, err := p.consume(RIGHT_BRACE, "Expect a '}' at the end of a block")
	if err != nil {
		return statements, err
	}
	return statements, nil
}

func (p *Parser[T]) printStatement() (Stmt[T], error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(SEMICOLON, "expect ';' at the end of a statement")
	if err != nil {
		return nil, err
	}
	return NewPrint(expr), nil
}

func (p *Parser[T]) expressionStatement() (Stmt[T], error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(SEMICOLON, "expect ';' at the end of a statement")
	if err != nil {
		return nil, err
	}
	return NewExpression(expr), nil
}

func (p *Parser[T]) varDeclaration() (Stmt[T], error) {
	name, err := p.consume(IDENTIFIER, "expect an identifier after var")
	if err != nil {
		return nil, err
	}
	var initializer Expr[T]
	if p.match(EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(SEMICOLON, "expect a ';' at the end of a declaration of a variable")
	if err != nil {
		return nil, err
	}
	return NewVar(name, initializer), nil
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
	return p.assignment()
}

func (p *Parser[T]) assignment() (Expr[T], error) {
	expr, err := p.equality()
	if p.match(EQUAL) {
		equals := p.previous()
		expr, ok := expr.(*Variable[T])
		if ok {
			value, err := p.assignment()
			if err != nil {
				return nil, err
			}
			return NewAssign(expr.name, value), nil
		}
		return nil, NewSyntaxError(equals.line, "Invalid assignment target.")
	}
	return expr, err
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
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return NewUnary(operator, right), nil
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
	} else if p.match(IDENTIFIER) {
		return NewVariable[T](p.previous()), nil
	}
	var msg string
	if p.current != 0 {
		msg = fmt.Sprintf("%q is not a valid expression", p.previous().ToString())
	} else {
		msg = fmt.Sprintf("%q is not a valid expression", p.peek().ToString())
	}
	return nil, NewSyntaxError(p.peek().line, msg)
}

func (p *Parser[T]) consume(tokenType TokenType, expectMessage string) (*Token, error) {
	if p.check(tokenType) {
		return p.next(), nil
	}
	return p.peek(), NewSyntaxError(p.peek().line, expectMessage)
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

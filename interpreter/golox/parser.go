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
	} else if p.match(IF) {
		return p.ifStatement()
	} else if p.match(WHILE) {
		return p.whileStatement()
	} else if p.match(FOR) {
		return p.forStatement()
	} else if p.match(LEFT_BRACE) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return NewBlock(statements), nil
	}
	return p.expressionStatement()
}

func (p *Parser[T]) forStatement() (Stmt[T], error) {
	_, err := p.consume(LEFT_PAREN, "Missing parenthesis before the clauses of the for statement")
	if err != nil {
		return nil, err
	}
	var initializer Stmt[T]
	if p.match(VAR) {
		initializer, err = p.varDeclaration()
	} else if !p.match(SEMICOLON) {
		initializer, err = p.expressionStatement()
	}
	if err != nil {
		return nil, err
	}
	var condition Expr[T]
	if !p.check(SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	p.consume(SEMICOLON, "expect a ';' at the end of a condition")
	var increment Expr[T]
	if !p.check(RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	p.consume(RIGHT_PAREN, "expect a ')' at the end of a condition")
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	if increment != nil {
		body = NewBlock([]Stmt[T]{
			body,
			NewExpression(increment),
		})
	}
	if condition != nil {
		body = NewWhile(condition, body)
	}
	if initializer != nil {
		body = NewBlock([]Stmt[T]{initializer, body})
	}
	return body, nil
}

func (p *Parser[T]) whileStatement() (Stmt[T], error) {
	_, err := p.consume(LEFT_PAREN, "Missing parenthesis before the condition of the while statement")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(RIGHT_PAREN, "Missing parenthesis after the condition of the while statement")
	if err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return NewWhile(condition, body), nil
}

func (p *Parser[T]) ifStatement() (Stmt[T], error) {
	_, err := p.consume(LEFT_PAREN, "Missing parenthesis before the condition of the if statement")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(RIGHT_PAREN, "Missing parenthesis after the condition of the if statement")
	if err != nil {
		return nil, err
	}
	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch Stmt[T]
	if p.match(ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return NewIf(condition, thenBranch, elseBranch), nil
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
	expr, err := p.logicalOr()
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

func (p *Parser[T]) logicalOr() (Expr[T], error) {
	leftExpr, err := p.logicalAnd()
	if err != nil {
		return nil, err
	}
	var rightExpr Expr[T]
	for {
		if !p.match(OR) {
			break
		}
		operator := p.previous()
		rightExpr, err = p.logicalAnd()
		if err != nil {
			return nil, err
		}
		leftExpr = NewLogical(leftExpr, operator, rightExpr)
	}
	return leftExpr, nil
}

func (p *Parser[T]) logicalAnd() (Expr[T], error) {
	leftExpr, err := p.equality()
	if err != nil {
		return nil, err
	}
	var rightExpr Expr[T]
	for {
		if !p.match(AND) {
			break
		}
		operator := p.previous()
		rightExpr, err = p.equality()
		if err != nil {
			return nil, err
		}
		leftExpr = NewLogical(leftExpr, operator, rightExpr)
	}
	return leftExpr, nil
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
		return p.call()
	}
}

func (p *Parser[T]) call() (Expr[T], error) {
	callee, err := p.primary()
	if err != nil {
		return nil, err
	}
	for {
		if !p.match(LEFT_PAREN) {
			break
		}
		paren := p.previous()
		var arguments []Expr[T]
		if !p.check(RIGHT_PAREN) {
			arguments, err = p.arguments()
			if err != nil {
				return nil, err
			}
		}
		_, err = p.consume(RIGHT_PAREN, "expect ')' after the arguments")
		if err != nil {
			return nil, err
		}
		callee = NewCall(callee, paren, arguments)
	}
	return callee, nil
}
func (p *Parser[T]) arguments() ([]Expr[T], error) {
	arguments := make([]Expr[T], 0, 10)
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	arguments = append(arguments, expr)
	for {
		if !p.match(COMMA) {
			break
		}
		if len(arguments) > 255 {
			return nil, NewSyntaxError(p.peek().line, "can't have more than 255 arguments")
		}
		expr, err = p.expression()
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, expr)
	}
	return arguments, nil
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

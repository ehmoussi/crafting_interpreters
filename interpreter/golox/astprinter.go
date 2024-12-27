package golox

import (
	"errors"
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (ap *AstPrinter) Print(statements []Stmt[any]) string {
	errorMsg := "the construction of the AST failed"
	var printStatements strings.Builder
	for _, statement := range statements {
		value, err := statement.accept(ap)
		if err != nil {
			return err.Error()
		}
		valueString, ok := value.(string)
		if !ok {
			return errorMsg
		}
		printStatements.WriteString(valueString + "\n")
	}
	return printStatements.String()
}

func (ap *AstPrinter) visitBlockStmt(stmt *Block[any]) (any, error) {
	var builder strings.Builder
	builder.WriteString("{\n")
	for _, statement := range stmt.statements {
		value, err := statement.accept(ap)
		if err != nil {
			return "", err
		}
		valueString, ok := value.(string)
		if !ok {
			return nil, errors.New("the construction of the AST failed")
		}
		builder.WriteString(valueString + "\n")
	}
	builder.WriteString("}")
	return builder.String(), nil
}

func (ap *AstPrinter) visitExpressionStmt(stmt *Expression[any]) (any, error) {
	return stmt.expression.accept(ap)
}

func (ap *AstPrinter) visitIfStmt(stmt *If[any]) (any, error) {
	var builder strings.Builder
	builder.WriteString("(")
	condition, err := ap.parenthesize("if", stmt.condition)
	if err != nil {
		return nil, err
	}
	builder.WriteString(condition)
	thenBranch, err := stmt.thenBranch.accept(ap)
	if err != nil {
		return nil, err
	}
	thenBranchString, ok := thenBranch.(string)
	if !ok {
		return nil, errors.New("the construction of the AST failed")
	}
	builder.WriteString(thenBranchString)
	if stmt.elseBranch != nil {
		elseBranch, err := stmt.elseBranch.accept(ap)
		if err != nil {
			return nil, err
		}
		elseBranchString, ok := elseBranch.(string)
		if !ok {
			return nil, errors.New("the construction of the AST failed")
		}
		builder.WriteString(elseBranchString)
	}
	builder.WriteString(")")
	return builder.String(), nil
}

func (ap *AstPrinter) visitPrintStmt(stmt *Print[any]) (any, error) {
	return ap.parenthesize("print", stmt.expression)
}

func (ap *AstPrinter) visitVarStmt(stmt *Var[any]) (any, error) {
	return ap.parenthesize("var "+stmt.name.lexeme, stmt.initializer)
}

func (ap *AstPrinter) visitAssignExpr(expr *Assign[any]) (any, error) {
	var builder strings.Builder
	builder.WriteString("(=")
	assignement, err := ap.parenthesize(expr.name.lexeme, expr.value)
	if err != nil {
		return "", err
	}
	builder.WriteString(assignement)
	builder.WriteString(")")
	return builder.String(), nil
}

func (ap *AstPrinter) visitBinaryExpr(expr *Binary[any]) (any, error) {
	return ap.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (ap *AstPrinter) visitGroupingExpr(expr *Grouping[any]) (any, error) {
	return ap.parenthesize("group", expr.expression)
}

func (ap *AstPrinter) visitLiteralExpr(expr *Literal[any]) (any, error) {
	if expr.value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.value), nil
}

func (ap *AstPrinter) visitLogicalExpr(expr *Logical[any]) (any, error) {
	return ap.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (ap *AstPrinter) visitUnaryExpr(expr *Unary[any]) (any, error) {
	return ap.parenthesize(expr.operator.lexeme, expr.right)
}

func (ap *AstPrinter) visitVariableExpr(expr *Variable[any]) (any, error) {
	return expr.name.lexeme, nil
}

func (ap *AstPrinter) parenthesize(name string, exprs ...Expr[any]) (string, error) {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		e, err := expr.accept(ap)
		if err != nil {
			return "", err
		}
		eString, ok := e.(string)
		if !ok {
			return "", errors.New("unexpected error: the expression is not a string")
		}
		builder.WriteString(eString)
	}
	builder.WriteString(")")
	return builder.String(), nil
}

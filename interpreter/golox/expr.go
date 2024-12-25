package golox

type Expr[T any] interface {
    accept(visitor Visitor[T]) T
}

type Binary[T any] struct {
    left Expr[T]
    operator *Token
    right Expr[T]
}

func NewBinary[T any](left Expr[T], operator *Token, right Expr[T]) *Binary[T] {
    return &Binary[T]{
        left: left,
        operator: operator,
        right: right,
    }
}

func (e *Binary[T]) accept(visitor Visitor[T]) T{
    return visitor.visitBinaryExpr(e)}

type Grouping[T any] struct {
    expression Expr[T]
}

func NewGrouping[T any](expression Expr[T]) *Grouping[T] {
    return &Grouping[T]{
        expression: expression,
    }
}

func (e *Grouping[T]) accept(visitor Visitor[T]) T{
    return visitor.visitGroupingExpr(e)}

type Literal[T any] struct {
    value any
}

func NewLiteral[T any](value any) *Literal[T] {
    return &Literal[T]{
        value: value,
    }
}

func (e *Literal[T]) accept(visitor Visitor[T]) T{
    return visitor.visitLiteralExpr(e)}

type Unary[T any] struct {
    operator *Token
    right Expr[T]
}

func NewUnary[T any](operator *Token, right Expr[T]) *Unary[T] {
    return &Unary[T]{
        operator: operator,
        right: right,
    }
}

func (e *Unary[T]) accept(visitor Visitor[T]) T{
    return visitor.visitUnaryExpr(e)}

type Visitor[T any] interface {
    visitBinaryExpr(expr *Binary[T]) T
    visitGroupingExpr(expr *Grouping[T]) T
    visitLiteralExpr(expr *Literal[T]) T
    visitUnaryExpr(expr *Unary[T]) T
}


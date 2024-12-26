package golox

type Expr[T any] interface {
    accept(visitor Visitor[T]) (T, error)
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

func (e *Binary[T]) accept(visitor Visitor[T]) (T, error){
    return visitor.visitBinaryExpr(e)}

type Grouping[T any] struct {
    expression Expr[T]
}

func NewGrouping[T any](expression Expr[T]) *Grouping[T] {
    return &Grouping[T]{
        expression: expression,
    }
}

func (e *Grouping[T]) accept(visitor Visitor[T]) (T, error){
    return visitor.visitGroupingExpr(e)}

type Literal[T any] struct {
    value any
}

func NewLiteral[T any](value any) *Literal[T] {
    return &Literal[T]{
        value: value,
    }
}

func (e *Literal[T]) accept(visitor Visitor[T]) (T, error){
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

func (e *Unary[T]) accept(visitor Visitor[T]) (T, error){
    return visitor.visitUnaryExpr(e)}

type Visitor[T any] interface {
    visitBinaryExpr(expr *Binary[T]) (T, error)
    visitGroupingExpr(expr *Grouping[T]) (T, error)
    visitLiteralExpr(expr *Literal[T]) (T, error)
    visitUnaryExpr(expr *Unary[T]) (T, error)
}


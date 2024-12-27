package golox

type Expr[T any] interface {
    accept(visitor ExprVisitor[T]) (T, error)
}

type Assign[T any] struct {
    name *Token
    value Expr[T]
}

func NewAssign[T any](name *Token, value Expr[T]) *Assign[T] {
    return &Assign[T]{
        name: name,
        value: value,
    }
}

func (e *Assign[T]) accept(visitor ExprVisitor[T]) (T, error){
    return visitor.visitAssignExpr(e)
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

func (e *Binary[T]) accept(visitor ExprVisitor[T]) (T, error){
    return visitor.visitBinaryExpr(e)
}

type Grouping[T any] struct {
    expression Expr[T]
}

func NewGrouping[T any](expression Expr[T]) *Grouping[T] {
    return &Grouping[T]{
        expression: expression,
    }
}

func (e *Grouping[T]) accept(visitor ExprVisitor[T]) (T, error){
    return visitor.visitGroupingExpr(e)
}

type Literal[T any] struct {
    value any
}

func NewLiteral[T any](value any) *Literal[T] {
    return &Literal[T]{
        value: value,
    }
}

func (e *Literal[T]) accept(visitor ExprVisitor[T]) (T, error){
    return visitor.visitLiteralExpr(e)
}

type Logical[T any] struct {
    left Expr[T]
    operator *Token
    right Expr[T]
}

func NewLogical[T any](left Expr[T], operator *Token, right Expr[T]) *Logical[T] {
    return &Logical[T]{
        left: left,
        operator: operator,
        right: right,
    }
}

func (e *Logical[T]) accept(visitor ExprVisitor[T]) (T, error){
    return visitor.visitLogicalExpr(e)
}

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

func (e *Unary[T]) accept(visitor ExprVisitor[T]) (T, error){
    return visitor.visitUnaryExpr(e)
}

type Variable[T any] struct {
    name *Token
}

func NewVariable[T any](name *Token) *Variable[T] {
    return &Variable[T]{
        name: name,
    }
}

func (e *Variable[T]) accept(visitor ExprVisitor[T]) (T, error){
    return visitor.visitVariableExpr(e)
}

type ExprVisitor[T any] interface {
    visitAssignExpr(expr *Assign[T]) (T, error)
    visitBinaryExpr(expr *Binary[T]) (T, error)
    visitGroupingExpr(expr *Grouping[T]) (T, error)
    visitLiteralExpr(expr *Literal[T]) (T, error)
    visitLogicalExpr(expr *Logical[T]) (T, error)
    visitUnaryExpr(expr *Unary[T]) (T, error)
    visitVariableExpr(expr *Variable[T]) (T, error)
}


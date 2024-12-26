package golox

type Stmt[T any] interface {
    accept(visitor StmtVisitor[T]) (T, error)
}

type Expression[T any] struct {
    expression Expr[T]
}

func NewExpression[T any](expression Expr[T]) *Expression[T] {
    return &Expression[T]{
        expression: expression,
    }
}

func (e *Expression[T]) accept(visitor StmtVisitor[T]) (T, error){
    return visitor.visitExpressionStmt(e)
}

type Print[T any] struct {
    expression Expr[T]
}

func NewPrint[T any](expression Expr[T]) *Print[T] {
    return &Print[T]{
        expression: expression,
    }
}

func (e *Print[T]) accept(visitor StmtVisitor[T]) (T, error){
    return visitor.visitPrintStmt(e)
}

type StmtVisitor[T any] interface {
    visitExpressionStmt(expr *Expression[T]) (T, error)
    visitPrintStmt(expr *Print[T]) (T, error)
}


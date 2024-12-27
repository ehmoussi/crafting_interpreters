package golox

type Stmt[T any] interface {
    accept(visitor StmtVisitor[T]) (T, error)
}

type Block[T any] struct {
    statements []Stmt[T]
}

func NewBlock[T any](statements []Stmt[T]) *Block[T] {
    return &Block[T]{
        statements: statements,
    }
}

func (e *Block[T]) accept(visitor StmtVisitor[T]) (T, error){
    return visitor.visitBlockStmt(e)
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

type Var[T any] struct {
    name *Token
    initializer Expr[T]
}

func NewVar[T any](name *Token, initializer Expr[T]) *Var[T] {
    return &Var[T]{
        name: name,
        initializer: initializer,
    }
}

func (e *Var[T]) accept(visitor StmtVisitor[T]) (T, error){
    return visitor.visitVarStmt(e)
}

type StmtVisitor[T any] interface {
    visitBlockStmt(expr *Block[T]) (T, error)
    visitExpressionStmt(expr *Expression[T]) (T, error)
    visitPrintStmt(expr *Print[T]) (T, error)
    visitVarStmt(expr *Var[T]) (T, error)
}


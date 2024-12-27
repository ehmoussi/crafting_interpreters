package golox

type GoLoxCallable interface {
	arity() int
	call(*Interpreter, []any) (any, error)
	String() string
}

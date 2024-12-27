package golox

import "fmt"

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment() *Environment {
	values := make(map[string]any, 0)
	return &Environment{enclosing: nil, values: values}
}

func NewEnvironmentWithEnclosing(enclosing *Environment) *Environment {
	values := make(map[string]any, 0)
	return &Environment{enclosing: enclosing, values: values}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) get(name *Token) (any, error) {
	value, ok := e.values[name.lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.get(name)
		}
		return nil, NewRuntimeError(
			name,
			fmt.Sprintf("Undefined variable '%s'", name.lexeme),
		)
	}
	return value, nil
}

func (e *Environment) assign(name *Token, value any) error {
	_, ok := e.values[name.lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.assign(name, value)
		}
		return NewRuntimeError(
			name,
			fmt.Sprintf("Undefined variable '%s'", name.lexeme),
		)
	}
	e.values[name.lexeme] = value
	return nil
}

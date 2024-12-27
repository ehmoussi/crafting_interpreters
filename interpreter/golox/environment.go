package golox

import "fmt"

type Environment struct {
	values map[string]any
}

func NewEnvironment() *Environment {
	values := make(map[string]any, 0)
	return &Environment{values: values}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) get(name *Token) (any, error) {
	value, ok := e.values[name.lexeme]
	if !ok {
		return nil, NewRuntimeError(
			name,
			fmt.Sprintf("Undefined variable '%s'", name.lexeme),
		)
	}
	return value, nil
}

func (e *Environment) assign(name *Token, value any) error {
	_, err := e.get(name)
	if err != nil {
		return err
	}
	e.values[name.lexeme] = value
	return nil
}

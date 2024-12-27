package golox

import (
	"fmt"
	"os"
	"time"
)

type Clock struct{}

func (c *Clock) arity() int { return 0 }

func (c *Clock) call(interp *Interpreter, args []any) (any, error) {
	currentTime := time.Now()
	return float64(currentTime.UnixMilli()) / 1000.0, nil
}

func (c *Clock) String() string {
	return "<native function>"
}

type ReadFile struct{}

func (c *ReadFile) arity() int { return 1 }

func (c *ReadFile) call(interp *Interpreter, args []any) (any, error) {
	path, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("the first argument must be a string: %s", args[0])
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (c *ReadFile) String() string {
	return "<native function>"
}

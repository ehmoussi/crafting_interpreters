package golox

import "time"

type Clock struct{}

func (c *Clock) arity() int { return 0 }

func (c *Clock) call(interp *Interpreter, args []any) (any, error) {
	currentTime := time.Now()
	return float64(currentTime.UnixMilli()) / 1000.0, nil
}

func (c *Clock) String() string {
	return "<native function>"
}

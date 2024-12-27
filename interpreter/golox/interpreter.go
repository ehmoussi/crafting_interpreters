package golox

import "fmt"

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{environment: NewEnvironment()}
}

func (interp *Interpreter) interpret(statements []Stmt[any], isRepl bool) ([]string, error) {
	capacity := 0
	if isRepl {
		capacity = 50
	}
	values := make([]string, 0, capacity)
	for _, stmt := range statements {
		value, err := interp.execute(stmt)
		if err != nil {
			return nil, err
		}
		if value != nil && isRepl {
			values = append(values, interp.stringify(value))
		}
	}
	return values, nil
}

func (interp *Interpreter) visitVarStmt(stmt *Var[any]) (any, error) {
	var value any
	if stmt.initializer != nil {
		var err error
		value, err = interp.evaluate(stmt.initializer)
		if err != nil {
			return nil, err
		}
	}
	interp.environment.define(stmt.name.lexeme, value)
	return value, nil
}

func (interp *Interpreter) visitBlockStmt(expr *Block[any]) (any, error) {
	values := make([]any, 0, 50)
	currentEnv := interp.environment
	interp.environment = NewEnvironmentWithEnclosing(currentEnv)
	for _, statement := range expr.statements {
		value, err := interp.execute(statement)
		if err != nil {
			interp.environment = currentEnv
			return nil, err
		}
		values = append(values, value)
	}
	interp.environment = currentEnv
	return values, nil
}

func (interp *Interpreter) visitExpressionStmt(stmt *Expression[any]) (any, error) {
	value, err := interp.evaluate(stmt.expression)
	return value, err
}

func (interp *Interpreter) visitIfStmt(stmt *If[any]) (any, error) {
	value, err := interp.evaluate(stmt.condition)
	if err != nil {
		return nil, err
	}
	if interp.isTruthy(value) {
		value, err := interp.execute(stmt.thenBranch)
		if err != nil {
			return nil, err
		}
		return value, nil
	} else if stmt.elseBranch != nil {
		value, err := interp.execute(stmt.elseBranch)
		if err != nil {
			return nil, err
		}
		return value, nil
	}
	return nil, nil
}

func (interp *Interpreter) visitPrintStmt(stmt *Print[any]) (any, error) {
	value, err := interp.evaluate(stmt.expression)
	if err == nil {
		fmt.Println(interp.stringify(value))
	}
	return nil, err
}

func (interp *Interpreter) execute(stmt Stmt[any]) (any, error) {
	value, err := stmt.accept(interp)
	return value, err
}

func (interp *Interpreter) visitLiteralExpr(expr *Literal[any]) (any, error) {
	return expr.value, nil
}

func (interp *Interpreter) visitLogicalExpr(expr *Logical[any]) (any, error) {
	leftValue, err := interp.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	if expr.operator.tokenType == AND && !interp.isTruthy(leftValue) {
		return leftValue, nil
	} else if expr.operator.tokenType == OR && interp.isTruthy(leftValue) {
		return leftValue, nil
	}
	return interp.evaluate(expr.right)
}

func (interp *Interpreter) visitGroupingExpr(expr *Grouping[any]) (any, error) {
	return interp.evaluate(expr.expression)
}

func (interp *Interpreter) visitVariableExpr(expr *Variable[any]) (any, error) {
	return interp.environment.get(expr.name)
}

func (interp *Interpreter) visitUnaryExpr(expr *Unary[any]) (any, error) {
	right, _ := interp.evaluate(expr.right)
	switch expr.operator.tokenType {
	case BANG:
		return !interp.isTruthy(right), nil
	case MINUS:
		rightValue, err := interp.checkNumberOperand(expr.operator, right)
		if err != nil {
			return nil, err
		}
		return -rightValue, nil
	}
	return nil, NewRuntimeError(expr.operator, "The operator is not valid for unary expression")
}

func (interp *Interpreter) visitBinaryExpr(expr *Binary[any]) (any, error) {
	left, _ := interp.evaluate(expr.left)
	right, _ := interp.evaluate(expr.right)
	switch expr.operator.tokenType {
	case MINUS:
		leftValue, rightValue, err := interp.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftValue - rightValue, nil
	case SLASH:
		leftValue, rightValue, err := interp.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftValue / rightValue, nil
	case STAR:
		leftValue, rightValue, err := interp.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftValue * rightValue, nil
	case GREATER:
		leftValue, rightValue, err := interp.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftValue > rightValue, nil
	case GREATER_EQUAL:
		leftValue, rightValue, err := interp.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftValue >= rightValue, nil
	case LESS:
		leftValue, rightValue, err := interp.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftValue < rightValue, nil
	case LESS_EQUAL:
		leftValue, rightValue, err := interp.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftValue <= rightValue, nil
	case PLUS:
		if leftValue, okLeft := left.(float64); okLeft {
			if rightValue, okRight := right.(float64); okRight {
				return leftValue + rightValue, nil
			}
		}
		if leftValue, okLeft := left.(string); okLeft {
			if rightValue, okRight := right.(string); okRight {
				return leftValue + rightValue, nil
			}
		}
		return nil, NewRuntimeError(expr.operator, "The operands must be two numbers or two strings")
	case BANG_EQUAL:
		return left != right, nil
	case EQUAL_EQUAL:
		return left == right, nil
	}
	return nil, NewRuntimeError(expr.operator, "The operator is no valid for binary expression")
}

func (interp *Interpreter) visitAssignExpr(expr *Assign[any]) (any, error) {
	value, err := interp.evaluate(expr.value)
	if err != nil {
		return nil, err
	}
	err = interp.environment.assign(expr.name, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (interp *Interpreter) evaluate(expr Expr[any]) (any, error) {
	return expr.accept(interp)
}

func (interp *Interpreter) isTruthy(value any) bool {
	switch value := value.(type) {
	case nil:
		return false
	case bool:
		return value
	default:
		return true
	}
}

func (interp *Interpreter) checkNumberOperands(operator *Token, left any, right any) (float64, float64, error) {
	leftValue, okLeft := left.(float64)
	rightValue, okRight := right.(float64)
	if !okLeft && !okRight {
		return 0.0, 0.0, NewRuntimeError(operator, "Operands must be numbers")
	} else if !okLeft {
		return 0.0, rightValue, NewRuntimeError(operator, "Left operand must be a number")
	} else if !okRight {
		return leftValue, 0.0, NewRuntimeError(operator, "Right operand must be a number")
	}
	return leftValue, rightValue, nil
}

func (interp *Interpreter) checkNumberOperand(operator *Token, operand any) (float64, error) {
	value, ok := operand.(float64)
	if !ok {
		return 0.0, NewRuntimeError(operator, "Operand must be a number")
	}
	return value, nil
}

func (interp *Interpreter) stringify(value any) string {
	switch value := value.(type) {
	case bool:
		return fmt.Sprintf("%t", value)
	case float64:
		return fmt.Sprintf("%f", value)
	default:
		return fmt.Sprintf("%s", value)
	}
}

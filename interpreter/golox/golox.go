package golox

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type GoLox struct {
	hadRuntimeError bool
}

func NewGoLox() *GoLox {
	return &GoLox{hadRuntimeError: false}
}

func (lox *GoLox) RunFile(path string) {
	interpreter := NewInterpreter()
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	lox.run(string(bytes), interpreter, false)
	if err != nil {
		fmt.Println(err)
		var syntaxErr *SyntaxError
		var runtimeErr *RuntimeError
		if errors.As(err, &syntaxErr) {
			os.Exit(65)
		} else if errors.As(err, &runtimeErr) {
			os.Exit(70)
		} else {
			os.Exit(1)
		}
	}
}

func (lox *GoLox) RunPrompt() {
	interpreter := NewInterpreter()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
		}
		err = lox.run(line, interpreter, true)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (lox *GoLox) run(source string, interpreter *Interpreter, isRepl bool) error {
	// Find the tokens
	scanner := NewScanner(source, 100)
	tokens, err := scanner.scanTokens()
	if err != nil {
		// fmt.Print(err)
		return err
	}
	// Parse the tokens
	parser := NewParser[any](len(tokens))
	parser.tokens = append(parser.tokens, tokens...)
	statements, err := parser.Parse()
	if err != nil {
		return err
	}
	// Print the AST
	// fmt.Println(NewAstPrinter().Print(statements))
	// Interpret the expression
	values, err := interpreter.interpret(statements, isRepl)
	if err != nil {
		return err
	}
	if len(values) > 0 {
		fmt.Println(strings.Join(values, "\n"))
	}
	return nil
}

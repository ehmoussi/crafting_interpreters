package golox

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
)

type GoLox struct {
	hadRuntimeError bool
}

func NewGoLox() *GoLox {
	return &GoLox{hadRuntimeError: false}
}

func (lox *GoLox) RunFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	lox.run(string(bytes))
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
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
		}
		err = lox.run(line)
		if err != nil {
			fmt.Print(err)
		}
	}
}

func (lox *GoLox) run(source string) error {
	// Find the tokens
	scanner := NewScanner(source, 100)
	tokens, err := scanner.scanTokens()
	if err != nil {
		fmt.Print(err)
	}
	// Parse the tokens
	parser := NewParser[any](len(tokens))
	parser.tokens = append(parser.tokens, tokens...)
	expr, err := parser.Parse()
	if err != nil || expr == nil {
		return err
	}
	// Print the AST
	fmt.Println(NewAstPrinter().Print(expr))
	// Interpret the expression
	value, err := NewInterpreter().evaluate(expr)
	if err != nil {
		return err
	}
	fmt.Println(value)
	return nil
}

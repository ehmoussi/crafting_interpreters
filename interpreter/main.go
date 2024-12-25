package main

import (
	"fmt"
	"golox/golox"
	"log"
	"os"
)

func main() {
	// Test AST
	expr := golox.NewBinary(
		golox.NewUnary(
			golox.NewToken(golox.MINUS, "-", nil, 1),
			golox.NewLiteral[string](123),
		),
		golox.NewToken(golox.STAR, "*", nil, 1),
		golox.NewGrouping(golox.NewLiteral[string](45.67)),
	)
	fmt.Println(golox.NewAstPrinter().Print(expr))
	// Run GoLox interpreter
	goLox := golox.NewGoLox()
	if len(os.Args) > 2 {
		log.Fatal("Usage: golox [script]")
	} else if len(os.Args) == 2 {
		goLox.RunFile(os.Args[1])
	} else {
		goLox.RunPrompt()
	}
}

package main

import (
	"log"
	"os"
)

func main() {
	lox := NewGoLox()
	if len(os.Args) > 2 {
		log.Fatal("Usage: golox [script]")
	} else if len(os.Args) == 2 {
		lox.runFile(os.Args[1])
	} else {
		lox.runPrompt()
	}
}

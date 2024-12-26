package main

import (
	"golox/golox"
	"log"
	"os"
)

func main() {
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

package main

import (
	"golox/golox"
	"log"
	"os"
)

func main() {
	goLox := golox.NewGoLox()
	if len(os.Args) > 2 {
		log.Fatal("Usage: golox [script]")
	} else if len(os.Args) == 2 {
		goLox.RunFile(os.Args[1])
	} else {
		goLox.RunPrompt()
	}
}

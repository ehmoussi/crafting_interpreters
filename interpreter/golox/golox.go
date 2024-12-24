package golox

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type GoLox struct {
}

func NewGoLox() *GoLox {
	return &GoLox{}
}

func (lox *GoLox) RunFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	lox.run(string(bytes))
}

func (lox *GoLox) RunPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
		}
		lox.run(line)
	}
}

func (lox *GoLox) run(source string) {
	scanner := NewScanner(source, 100)
	tokens := scanner.scanTokens()
	for _, token := range tokens {
		fmt.Println(token.ToString())
	}
}

func error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
}

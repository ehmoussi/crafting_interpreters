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
		err = lox.run(line)
		fmt.Println(err)
	}
}

func ReportError(line int, message string) {
	report(line, "", message)
}

func (lox *GoLox) run(source string) error {
	// Find the tokens
	scanner := NewScanner(source, 100)
	tokens := scanner.scanTokens()
	// for _, token := range tokens {
	// 	fmt.Println(token.ToString())
	// }
	// Parse the tokens
	parser := NewParser[string](len(tokens))
	parser.tokens = append(parser.tokens, tokens...)
	expr, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse: %q", source)
	}
	fmt.Println(NewAstPrinter().Print(expr))
	return nil
}

func report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
}

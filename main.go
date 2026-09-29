package main

import (
	"fmt"
	"os"

	"github.com/Moses-Bejon/python-lexer/pythonlexer"
)

func main() {
	if err := run(); err != nil {
		if _, writeErr := fmt.Fprintln(os.Stderr, err); writeErr != nil {
			os.Exit(2)
		}
		os.Exit(1)
	}
}

func run() error {
	path := "lexer.py"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	source, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	file := &pythonlexer.SourceFile{Source: []rune(string(source))}
	tokens, err := pythonlexer.NewLexer(file).Tokenise()
	for _, token := range tokens {
		if _, writeErr := fmt.Printf("%-18s %q\n", token.Type, string(token.Lexeme)); writeErr != nil {
			return fmt.Errorf("write token output: %w", writeErr)
		}
	}
	if err != nil {
		return fmt.Errorf("tokenise %s: %w", path, err)
	}
	return nil
}

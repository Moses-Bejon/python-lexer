# python-lexer

A Go implementation of a Python lexer.

This has been moved here from the work-in-progress private repository to which it was contributed, panpiler.

## Run

Run the lexer against the included Python sample:

```sh
go run .
```

Pass a different Python file as the first argument:

```sh
go run . path/to/file.py
```

The lexer implementation is in `pythonlexer/`; `main.go` is a small command-line
runner that prints each token type and lexeme.

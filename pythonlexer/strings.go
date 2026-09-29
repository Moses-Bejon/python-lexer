package pythonlexer

import (
	"errors"
	"slices"
)

var (
	ErrExpectedEscapeSlash       = errors.New("expected '\\' before an escape sequence")
	ErrExpectedHexEscapeDigit    = errors.New("expected at least one digit for a hex escape sequence")
	ErrExpectedUnicodeEscape4    = errors.New("expected four hexadecimal digits after \\u")
	ErrExpectedUnicodeEscape8    = errors.New("expected eight hexadecimal digits after \\U")
	ErrInvalidEscapeSequence     = errors.New("invalid escape sequence")
	ErrUnterminatedStringLiteral = errors.New("syntax error: unterminated string literal")
	ErrInvalidFstringEscape      = errors.New("\\ without an escape sequence in an f-string literal")
	ErrMultipleFormatSpecs       = errors.New("cannot have multiple format specs in a replacement field")
	ErrUnescapedFstringBrace     = errors.New("unescaped } found in f-string literal")
	ErrNoStringOpeningQuotes     = errors.New("no string opening quotes")
	ErrNoFstringPrefix           = errors.New("no f-string prefix")
	ErrNoTstringPrefix           = errors.New("no t-string prefix")
)

func (l *Lexer) escapeSequence(raw bool) error {
	backtrackTo := l.current

	if l.expected("\\") != nil {
		return ErrExpectedEscapeSlash
	}

	if raw {
		if l.expected("\"") != nil {

			// deliberately uncaught error:
			l.expected("'")
		}
		return nil
	}

	switch l.file.Source[l.current] {
	case '\n', '\\', '\'', '"', 'a', 'b', 'f', 'n', 'r', 't', 'v':
		l.current++
		return nil

	case '0', '1', '2', '3', '4', '5', '6', '7':

		for range 3 {
			l.current++

			if !slices.Contains([]rune{'0', '1', '2', '3', '4', '5', '6', '7'}, l.file.Source[l.current]) {
				break
			}
		}
		return nil

	case 'x':

		atLeastOneDigit := false

		for range 2 {
			l.current++

			if l.hexDigit() != nil {
				break
			}
			atLeastOneDigit = true
		}

		if !atLeastOneDigit {
			l.current = backtrackTo

			return ErrExpectedHexEscapeDigit
		}

		return nil

	case 'N':

		l.current++

		e := l.expected("{")

		if e != nil {
			l.current = backtrackTo

			return e
		}

		for l.file.Source[l.current] != '}' {
			l.current++
		}

		e = l.expected("}")

		if e != nil {
			l.current = backtrackTo

			return e
		}

		return nil

	case 'u':
		l.current++

		for range 4 {
			if l.hexDigit() != nil {
				l.current = backtrackTo

				return ErrExpectedUnicodeEscape4
			}
		}

		return nil

	case 'U':
		l.current++

		for range 8 {
			if l.hexDigit() != nil {
				l.current = backtrackTo

				return ErrExpectedUnicodeEscape8
			}
		}

		return nil
	}

	l.current = backtrackTo
	return ErrInvalidEscapeSequence
}

func (l *Lexer) stringUsing(terminator string, stopAt []rune) error {
	for l.current < len(l.file.Source) {

		l.goUntil(stopAt)

		// skip a character if it is escaped
		if l.file.Source[l.current] == '\\' {
			// 2 as it skips the \, then skips the next character
			l.current += 2
		} else if l.file.Source[l.current] == '\n' {
			return ErrUnterminatedStringLiteral
		}

		if l.expected(terminator) == nil {
			return nil
		}
	}

	return ErrUnterminatedStringLiteral
}

func (l *Lexer) ftstringUsing(

	terminator string,
	stopAt []rune,
	startToken TokenType,
	middleToken TokenType,
	endToken TokenType,

	raw bool,

) ([]Token, error) {

	tokens := []Token{l.generateToken(startToken)}
	insideFormatSpec := []bool{false}

	l.start = l.current

	for l.current < len(l.file.Source) {

		l.goUntil(stopAt)

		// skip a character if it is escaped
		if l.file.Source[l.current] == '\\' {
			if l.escapeSequence(raw) != nil {
				return tokens, ErrInvalidFstringEscape
			}
		} else if l.file.Source[l.current] == '\n' {
			return tokens, ErrUnterminatedStringLiteral
		} else if l.file.Source[l.current] == '{' {
			if l.file.Source[l.current+1] == '{' {
				l.current += 2
			} else {

				if l.start != l.current {
					tokens = append(tokens, l.generateToken(middleToken))
					l.start = l.current
				}

				l.current++
				tokens = append(tokens, l.generateToken(LBRACE))

				insideFormatSpec = append(insideFormatSpec, false)

				// whitespace between open brace and expression is ignored
				l.removePythonWhitespace()

				replacementField := NewLexer(l.file)

				replacementField.InitialiseLexer()

				replacementField.current = l.current
				replacementField.braceDepth = 1

				for replacementField.current < len(replacementField.file.Source) {
					replacementField.start = replacementField.current

					e := replacementField.nextTokens()
					if e != nil {
						return tokens, e
					}

					if replacementField.braceDepth <= 0 {

						insideFormatSpec = insideFormatSpec[:len(insideFormatSpec)-1]
						break
					}

					if replacementField.tokens[len(replacementField.tokens)-1].Type == COLON &&
						replacementField.parenDepth <= 0 &&
						replacementField.braceDepth <= 1 &&
						replacementField.squareBracketDepth <= 0 {

						if insideFormatSpec[len(insideFormatSpec)-1] {
							return tokens, ErrMultipleFormatSpecs
						}

						insideFormatSpec[len(insideFormatSpec)-1] = true

						break
					}
				}

				tokens = append(tokens, replacementField.tokens...)

				l.current = replacementField.current
				l.start = l.current
			}
		} else if l.file.Source[l.current] == '}' {

			if l.file.Source[l.current+1] == '}' {
				l.current += 2
			} else {

				if !insideFormatSpec[len(insideFormatSpec)-1] {
					l.current++
					return tokens, ErrUnescapedFstringBrace
				}

				tokens = append(tokens, l.generateToken(middleToken))
				l.start = l.current
				l.current++

				tokens = append(tokens, l.generateToken(RBRACE))

				l.start = l.current

				insideFormatSpec = insideFormatSpec[:len(insideFormatSpec)-1]
			}
		}

		place := l.current

		if l.expected(terminator) == nil {

			if l.start != place {
				tokens = append(tokens, Token{
					Type:   middleToken,
					Offset: place,
					Source: l.file,
					Lexeme: l.file.Source[l.start:place],
				})
				l.start = place
			}

			tokens = append(tokens, l.generateToken(endToken))

			return tokens, nil
		}
	}

	return tokens, ErrUnterminatedStringLiteral
}

func (l *Lexer) pystring() error {

	backtrackTo := l.current

	// deliberately not caught errors
	l.expectedIgnoreCase("r")
	l.expectedIgnoreCase("b")

	if l.expected("\"\"\"") == nil {

		e := l.stringUsing("\"\"\"", []rune{'"', '\\'})

		if e != nil {
			l.current = backtrackTo
			return e
		}

	} else if l.expected("'''") == nil {

		e := l.stringUsing("'''", []rune{'\'', '\\'})

		if e != nil {
			l.current = backtrackTo
			return e
		}

	} else if l.expected("'") == nil {

		e := l.stringUsing("'", []rune{'\'', '\\', '\n'})

		if e != nil {
			l.current = backtrackTo
			return e
		}

	} else if l.expected("\"") == nil {

		e := l.stringUsing("\"", []rune{'"', '\\', '\n'})

		if e != nil {
			l.current = backtrackTo
			return e
		}

	} else {
		l.current = backtrackTo
		return ErrNoStringOpeningQuotes
	}

	return nil
}

func (l *Lexer) fstring() error {

	backtrackTo := l.current

	var tokens []Token
	var e error

	raw := l.expectedIgnoreCase("fr") == nil || l.expectedIgnoreCase("rf") == nil

	if !raw && l.expectedIgnoreCase("f") != nil {
		return ErrNoFstringPrefix
	}

	if l.expected("\"\"\"") == nil {

		tokens, e = l.ftstringUsing(
			"\"\"\"",
			[]rune{'"', '\\', '{', '}'},
			FSTRING_START,
			FSTRING_MIDDLE,
			FSTRING_END,
			raw,
		)

		if e != nil {

			l.start = backtrackTo
			l.current = backtrackTo

			return e
		}

	} else if l.expected("'''") == nil {

		tokens, e = l.ftstringUsing(
			"'''",
			[]rune{'\'', '\\', '{', '}'},
			FSTRING_START,
			FSTRING_MIDDLE,
			FSTRING_END,
			raw,
		)

		if e != nil {
			l.start = backtrackTo
			l.current = backtrackTo

			return e
		}

	} else if l.expected("'") == nil {

		tokens, e = l.ftstringUsing(
			"'",
			[]rune{'\'', '\\', '\n', '{', '}'},
			FSTRING_START,
			FSTRING_MIDDLE,
			FSTRING_END,
			raw,
		)

		if e != nil {
			l.start = backtrackTo
			l.current = backtrackTo

			return e
		}

	} else if l.expected("\"") == nil {

		tokens, e = l.ftstringUsing(
			"\"",
			[]rune{'"', '\\', '\n', '{', '}'},
			FSTRING_START,
			FSTRING_MIDDLE,
			FSTRING_END,
			raw,
		)

		if e != nil {
			l.start = backtrackTo
			l.current = backtrackTo

			return e
		}

	} else {

		l.start = backtrackTo
		l.current = backtrackTo

		return ErrNoFstringPrefix
	}

	l.tokens = append(l.tokens, tokens...)
	return nil
}

func (l *Lexer) tstring() error {

	backtrackTo := l.current

	var tokens []Token
	var e error

	raw := l.expectedIgnoreCase("tr") == nil || l.expectedIgnoreCase("rt") == nil

	if !raw && l.expectedIgnoreCase("t") != nil {
		return ErrNoTstringPrefix
	}

	if l.expected("\"\"\"") == nil {

		tokens, e = l.ftstringUsing(
			"\"\"\"",
			[]rune{'"', '\\', '{', '}'},
			TSTRING_START,
			TSTRING_MIDDLE,
			TSTRING_END,
			raw,
		)

		if e != nil {

			l.start = backtrackTo
			l.current = backtrackTo

			return e
		}

	} else if l.expected("'''") == nil {

		tokens, e = l.ftstringUsing(
			"'''",
			[]rune{'\'', '\\', '{', '}'},
			TSTRING_START,
			TSTRING_MIDDLE,
			TSTRING_END,
			raw,
		)

		if e != nil {
			l.start = backtrackTo
			l.current = backtrackTo

			return e
		}

	} else if l.expected("'") == nil {

		tokens, e = l.ftstringUsing(
			"'",
			[]rune{'\'', '\\', '\n', '{', '}'},
			TSTRING_START,
			TSTRING_MIDDLE,
			TSTRING_END,
			raw,
		)

		if e != nil {
			l.start = backtrackTo
			l.current = backtrackTo

			return e
		}

	} else if l.expected("\"") == nil {

		tokens, e = l.ftstringUsing(
			"\"",
			[]rune{'"', '\\', '\n', '{', '}'},
			TSTRING_START,
			TSTRING_MIDDLE,
			TSTRING_END,
			raw,
		)

		if e != nil {
			l.start = backtrackTo
			l.current = backtrackTo

			return e
		}

	} else {

		l.start = backtrackTo
		l.current = backtrackTo

		return ErrNoTstringPrefix
	}

	l.tokens = append(l.tokens, tokens...)
	return nil
}

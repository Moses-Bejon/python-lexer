package pythonlexer

import (
	"errors"
	"slices"
)

var (
	ErrUnbalancedBraces         = errors.New("unbalanced braces")
	ErrUnbalancedParens         = errors.New("unbalanced parentheses")
	ErrUnbalancedSquareBrackets = errors.New("unbalanced square brackets")
	ErrIndentation              = errors.New("indentation error: unindent does not match any outer indentation level")
	ErrUnrecognizedToken        = errors.New("unrecognized token")
)

type Lexer struct {
	file *SourceFile

	start   int
	current int

	indents []int

	squareBracketDepth int
	parenDepth         int
	braceDepth         int

	tokens []Token
}

func NewLexer(source *SourceFile) *Lexer {
	return &Lexer{
		file: source,
	}
}

func (l *Lexer) InitialiseLexer() {
	l.start = 0
	l.current = 0

	l.indents = []int{0}

	l.squareBracketDepth = 0
	l.parenDepth = 0
	l.braceDepth = 0

	l.tokens = []Token{}
}

func (l *Lexer) Tokenise() ([]Token, error) {

	l.InitialiseLexer()

	e := l.newLogicalLine()
	if e != nil {
		return l.tokens, e
	}

	for l.current < len(l.file.Source) {
		l.start = l.current

		e := l.nextTokens()

		if e != nil {
			return l.tokens, e
		}
	}

	l.tokens = append(l.tokens, Token{
		Type:   ENDMARKER,
		Offset: l.current,
		Source: l.file,
		Lexeme: []rune(""),
	})

	return l.tokens, nil

}

func (l *Lexer) generateToken(tokenType TokenType) Token {
	return Token{
		Type:   tokenType,
		Offset: l.current,
		Source: l.file,
		Lexeme: l.file.Source[l.start:l.current],
	}
}

func (l *Lexer) newLogicalLine() error {

	if l.current >= len(l.file.Source) {
		return nil
	}

	l.start = l.current

	indentMarker := l.file.Source[l.current]

	// according to the python spec
	// "A formfeed character may be present at the start of the line;
	// it will be ignored for the indentation calculations"
	if indentMarker == '\f' {
		l.current++

		if l.current >= len(l.file.Source) {
			return nil
		}

		indentMarker = l.file.Source[l.current]
	}

	currentIndent := 0

	if indentMarker == ' ' || indentMarker == '\t' {

		for l.current < len(l.file.Source) && l.file.Source[l.current] == indentMarker {
			l.current++
			currentIndent++
		}

	}

	// ignore indentation if last line of file
	if l.current >= len(l.file.Source) {
		return nil
	}

	if l.file.Source[l.current] == '#' {
		l.start = l.current

		l.goUntil([]rune("\n"))

		l.tokens = append(l.tokens, l.generateToken(COMMENT))
		l.start = l.current
	}

	// if line is empty, ignore line and process next one
	if l.file.Source[l.current] == '\n' {
		l.start = l.current
		l.current++
		l.tokens = append(l.tokens, l.generateToken(NL))

		err := l.newLogicalLine()
		if err != nil {
			return err
		}

		return nil
	}

	stackTop := len(l.indents) - 1

	if currentIndent > l.indents[stackTop] {
		l.indents = append(l.indents, currentIndent)
		l.tokens = append(l.tokens, l.generateToken(INDENT))
	} else {

		// currentIndent >= 0
		// l.indents[0] = 0
		// so it follows this terminates before stackTop < 0
		for currentIndent < l.indents[stackTop] {
			stackTop--
			l.tokens = append(l.tokens, l.generateToken(DEDENT))
		}

		if currentIndent != l.indents[stackTop] {
			return ErrIndentation
		}

		l.indents = l.indents[:stackTop+1]
	}

	return nil
}

func (l *Lexer) removePythonWhitespace() {
	for l.current < len(l.file.Source) && slices.Contains([]rune{'\t', ' ', '\f'}, l.file.Source[l.current]) {
		l.current++
	}
}

func (l *Lexer) nextTokens() error {

	tokenised := false

	switch l.file.Source[l.current] {
	case '\n':
		if l.expected("\n") == nil {
			tokenised = true
			if l.squareBracketDepth > 0 || l.parenDepth > 0 || l.braceDepth > 0 {
				l.tokens = append(l.tokens, l.generateToken(NL))
			} else {
				l.tokens = append(l.tokens, l.generateToken(NEWLINE))

				e := l.newLogicalLine()

				if e != nil {
					return e
				}
			}
		}

	case '\\':
		if l.expected("\\\n") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(NL))
		}

	case '#':
		if l.expected("#") == nil {
			tokenised = true
			l.goUntil([]rune("\n"))

			l.tokens = append(l.tokens, l.generateToken(COMMENT))

		}

	case 'r', 'R':

		if l.pystring() == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(STRING))
		} else if l.fstring() == nil {
			tokenised = true
		} else if l.tstring() == nil {
			tokenised = true
		} else if l.name() == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(NAME))
		}

	case 'b', 'B':
		if l.pystring() == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(STRING))
		} else if l.name() == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(NAME))
		}

	case '"', '\'':
		if l.pystring() == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(STRING))
		}

	case 't', 'T':
		if l.tstring() == nil {
			tokenised = true
		} else if l.name() == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(NAME))
		}

	case 'f', 'F':
		if l.fstring() == nil {
			tokenised = true
		} else if l.name() == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(NAME))
		}

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if l.numeric() == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(NUMBER))
		}

	case '}':
		if l.expected("}") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(RBRACE))
			l.braceDepth--

			if l.braceDepth < 0 {
				return ErrUnbalancedBraces
			}
		}

	case '{':
		if l.expected("{") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(LBRACE))
			l.braceDepth++
		}

	case '(':
		if l.expected("(") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(LPAR))
			l.parenDepth++

			if l.parenDepth < 0 {
				return ErrUnbalancedParens
			}
		}

	case ')':
		if l.expected(")") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(RPAR))
			l.parenDepth--
		}

	case '[':
		if l.expected("[") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(LSQB))
			l.squareBracketDepth++
		}

	case ']':
		if l.expected("]") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(RSQB))
			l.squareBracketDepth--

			if l.braceDepth < 0 {
				return ErrUnbalancedSquareBrackets
			}
		}

	case ':':
		if l.expected(":=") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(COLONEQUAL))
		} else if l.expected(":") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(COLON))
		}

	case '.':
		if l.numeric() == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(NUMBER))
		} else if l.expected("...") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(ELLIPSIS))
		} else if l.expected(".") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(DOT))
		}

	case '-':
		if l.expected("->") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(RARROW))
		} else if l.expected("-=") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(MINEQUAL))
		} else if l.expected("-") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(MINUS))
		}

	case '@':
		if l.expected("@=") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(ATEQUAL))
		} else if l.expected("@") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(AT))
		}

	case '/':
		if l.expected("//=") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(DOUBLESLASHEQUAL))
		} else if l.expected("//") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(DOUBLESLASH))
		} else if l.expected("/=") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(SLASHEQUAL))
		} else if l.expected("/") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(SLASH))
		}

	case '>':
		if l.expected(">>=") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(RIGHTSHIFTEQUAL))
		} else if l.expected(">>") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(RIGHTSHIFT))
		} else if l.expected(">=") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(GREATEREQUAL))
		} else if l.expected(">") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(GREATER))
		}

	case '<':
		if l.expected("<<=") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(LEFTSHIFTEQUAL))
		} else if l.expected("<<") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(LEFTSHIFT))
		} else if l.expected("<=") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(LESSEQUAL))
		} else if l.expected("<") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(LESS))
		}

	case '^':
		if l.expected("^=") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(CIRCUMFLEXEQUAL))
		} else if l.expected("^") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(CIRCUMFLEX))
		}

	case '|':
		if l.expected("|=") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(VBAREQUAL))
		} else if l.expected("|") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(VBAR))
		}

	case '&':
		if l.expected("&=") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(AMPEREQUAL))
		} else if l.expected("&") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(AMPER))
		}

	case '%':
		if l.expected("%=") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(PERCENTEQUAL))
		} else if l.expected("%") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(PERCENT))
		}

	case '*':
		if l.expected("**") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(DOUBLESTAR))
		} else if l.expected("*=") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(STAREQUAL))
		} else if l.expected("*") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(STAR))
		}

	case '~':
		if l.expected("~") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(TILDE))
		}

	case '!':
		if l.expected("!=") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(NOTEQUAL))
		} else if l.expected("!") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(EXCLAMATION))
		}

	case '=':
		if l.expected("==") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(EQEQUAL))
		} else if l.expected("=") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(EQUAL))
		}

	case '+':
		if l.expected("+=") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(PLUSEQUAL))
		} else if l.expected("+") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(PLUS))
		}

	case ';':
		if l.expected(";") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(SEMI))
		}

	case ',':
		if l.expected(",") == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(COMMA))
		}

	default:
		if l.name() == nil {
			tokenised = true
			l.tokens = append(l.tokens, l.generateToken(NAME))
		}

	}

	if !tokenised {
		l.current++
		l.tokens = append(l.tokens, l.generateToken(ILLEGAL))
		return ErrUnrecognizedToken
	}

	l.removePythonWhitespace()

	return nil
}

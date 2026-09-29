package pythonlexer

// Enum of token types for the lexer

const (
	NAME           TokenType = iota // Token value that indicates an identifier or keyword.
	NUMBER                          // Token value that indicates a numeric literal.
	STRING                          // Token value that indicates a string or byte literal, excluding formatted string literals. The token string is not interpreted: it includes the surrounding quotation marks and the prefix (if given); backslashes are included literally, without processing escape sequences.
	COMMENT                         // Token value used to indicate a comment.
	NEWLINE                         // Token value that indicates the end of a logical line.
	NL                              // Token value used to indicate a non-terminating newline. NL tokens are generated when a logical line of code is continued over multiple physical lines.
	INDENT                          // Token value used at the beginning of a logical line to indicate the start of an indented block.
	DEDENT                          // Token value used at the beginning of a logical line to indicate the end of an indented block.
	FSTRING_START                   // Token value used to indicate the beginning of an f-string literal.
	FSTRING_MIDDLE                  // Token value used for literal text inside an f-string literal, including format specifications.
	FSTRING_END                     // Token value used to indicate the end of a f-string.
	TSTRING_START                   // Token value used to indicate the beginning of a template string literal.
	TSTRING_MIDDLE                  // Token value used for literal text inside a template string literal including format specifications.
	TSTRING_END                     // Token value used to indicate the end of a template string literal.
	ENDMARKER                       // Token value that indicates the end of input. Used in top-level grammar rules.
	ILLEGAL                         // Unknown token.

	LPAR             // "("
	RPAR             // ")"
	LSQB             // "["
	RSQB             // "]"
	COLON            // ":"
	COMMA            // ","
	SEMI             // ";"
	PLUS             // "+"
	MINUS            // "-"
	STAR             // "*"
	SLASH            // "/"
	VBAR             // "|"
	AMPER            // "&"
	LESS             // "<"
	GREATER          // ">"
	EQUAL            // "="
	DOT              // "."
	PERCENT          // "%"
	LBRACE           // "{"
	RBRACE           // "}"
	EQEQUAL          // "=="
	NOTEQUAL         // "!="
	LESSEQUAL        // "<="
	GREATEREQUAL     // ">="
	TILDE            // "~"
	CIRCUMFLEX       // "^"
	LEFTSHIFT        // "<<"
	RIGHTSHIFT       // ">>"
	DOUBLESTAR       // "**"
	PLUSEQUAL        // "+="
	MINEQUAL         // "-="
	STAREQUAL        // "*="
	SLASHEQUAL       // "/="
	PERCENTEQUAL     // "%="
	AMPEREQUAL       // "&="
	VBAREQUAL        // "|="
	CIRCUMFLEXEQUAL  // "^="
	LEFTSHIFTEQUAL   // "<<="
	RIGHTSHIFTEQUAL  // ">>="
	DOUBLESLASH      // "//"
	DOUBLESLASHEQUAL // "//="
	AT               // "@"
	ATEQUAL          // "@="
	RARROW           // "->"
	ELLIPSIS         // "..."
	COLONEQUAL       // ":="
	EXCLAMATION      // "!"

)

type TokenType int

type SourceFile struct {
	Source []rune
}

type Token struct {
	Type   TokenType
	Offset int
	Source *SourceFile
	Lexeme []rune
}

var tokenTypeNames = [...]string{
	"NAME",
	"NUMBER",
	"STRING",
	"COMMENT",
	"NEWLINE",
	"NL",
	"INDENT",
	"DEDENT",
	"FSTRING_START",
	"FSTRING_MIDDLE",
	"FSTRING_END",
	"TSTRING_START",
	"TSTRING_MIDDLE",
	"TSTRING_END",
	"ENDMARKER",
	"ILLEGAL",
	"LPAR",
	"RPAR",
	"LSQB",
	"RSQB",
	"COLON",
	"COMMA",
	"SEMI",
	"PLUS",
	"MINUS",
	"STAR",
	"SLASH",
	"VBAR",
	"AMPER",
	"LESS",
	"GREATER",
	"EQUAL",
	"DOT",
	"PERCENT",
	"LBRACE",
	"RBRACE",
	"EQEQUAL",
	"NOTEQUAL",
	"LESSEQUAL",
	"GREATEREQUAL",
	"TILDE",
	"CIRCUMFLEX",
	"LEFTSHIFT",
	"RIGHTSHIFT",
	"DOUBLESTAR",
	"PLUSEQUAL",
	"MINEQUAL",
	"STAREQUAL",
	"SLASHEQUAL",
	"PERCENTEQUAL",
	"AMPEREQUAL",
	"VBAREQUAL",
	"CIRCUMFLEXEQUAL",
	"LEFTSHIFTEQUAL",
	"RIGHTSHIFTEQUAL",
	"DOUBLESLASH",
	"DOUBLESLASHEQUAL",
	"AT",
	"ATEQUAL",
	"RARROW",
	"ELLIPSIS",
	"COLONEQUAL",
	"EXCLAMATION",
}

func (t TokenType) String() string {
	if t < 0 || int(t) >= len(tokenTypeNames) {
		return "UNKNOWN"
	}
	return tokenTypeNames[t]
}

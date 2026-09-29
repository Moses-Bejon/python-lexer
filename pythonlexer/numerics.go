package pythonlexer

import (
	"errors"
)

var (
	ErrInvalidHexDigit = errors.New("invalid hex digit")
	ErrInvalidDigit    = errors.New("invalid digit")
	ErrInvalidOctDigit = errors.New("invalid oct digit")
	ErrInvalidBinDigit = errors.New("invalid bin digit")
	ErrInvalidExponent = errors.New("invalid exponent")
)

func (l *Lexer) hexDigit() error {
	r := l.file.Source[l.current]

	if ('0' <= r && r <= '9') ||
		('a' <= r && r <= 'f') ||
		('A' <= r && r <= 'F') {
		l.current++
		return nil
	}

	return ErrInvalidHexDigit
}

func (l *Lexer) digit() error {
	r := l.file.Source[l.current]

	if '0' <= r && r <= '9' {
		l.current++
		return nil
	}

	return ErrInvalidDigit
}

func (l *Lexer) octDigit() error {
	r := l.file.Source[l.current]

	if '0' <= r && r <= '7' {
		l.current++
		return nil
	}

	return ErrInvalidOctDigit
}

func (l *Lexer) binDigit() error {
	r := l.file.Source[l.current]

	if r == '0' || r == '1' {
		l.current++
		return nil
	}

	return ErrInvalidBinDigit
}

func (l *Lexer) digitPart() error {
	if l.digit() != nil {
		return ErrInvalidDigit
	}

	for l.current < len(l.file.Source) {
		l.expected("_")
		if l.digit() != nil {
			return nil
		}
	}

	return nil
}

func (l *Lexer) exponent() error {
	backtrackTo := l.current

	if l.expectedIgnoreCase("e") != nil {
		return ErrInvalidExponent
	}

	if l.expected("+") != nil {
		l.expected("-")
	}

	if l.digitPart() != nil {
		l.current = backtrackTo
		return ErrInvalidDigit
	}

	return nil
}

func (l *Lexer) float() error {
	backtrackTo := l.current

	if l.expected(".") == nil {
		if l.digitPart() != nil {
			l.current = backtrackTo
			return ErrInvalidDigit
		}

		l.exponent()

		return nil
	}

	if l.digitPart() != nil {
		l.current = backtrackTo
		return ErrInvalidDigit
	}

	if l.expected(".") == nil {
		l.digitPart()
		l.exponent()
		return nil
	}

	if l.exponent() != nil {
		l.current = backtrackTo
		return ErrInvalidExponent
	}

	return nil
}

func (l *Lexer) numeric() error {

	backtrackTo := l.current

	// floats and imaginaries

	if l.float() == nil {
		l.expectedIgnoreCase("j")
		return nil
	}

	if l.digitPart() == nil {
		if l.expectedIgnoreCase("j") == nil {
			return nil
		}
		l.current = backtrackTo
	}

	// ints

	bin := l.expectedIgnoreCase("0b") == nil
	oct := l.expectedIgnoreCase("0o") == nil
	hex := l.expectedIgnoreCase("0x") == nil

	var parser func() error = nil

	if bin {
		parser = l.binDigit
	} else if oct {
		parser = l.octDigit
	} else if hex {
		parser = l.hexDigit
	}

	if parser != nil {
		l.expected("_")
		if parser() != nil {
			l.current = backtrackTo
			return ErrInvalidDigit
		}

		for l.current < len(l.file.Source) {
			l.expected("_")

			if parser() != nil {
				return nil
			}
		}
		return nil
	}

	zeroInteger := l.expected("0") == nil

	if zeroInteger {
		for l.current < len(l.file.Source) {
			l.expected("_")

			if l.expected("0") != nil {
				return nil
			}
		}
		return nil
	}

	if l.digitPart() != nil {
		l.current = backtrackTo
		return ErrInvalidDigit
	}

	return nil
}

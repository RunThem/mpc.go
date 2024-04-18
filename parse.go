package mpc

import (
	"bytes"

	"github.com/RunThem/u"
)

// func New(code string) *Mpc {
//
// }

/*
 * token
 */
const (
	tok_IDENT = "IDENT"
	tok_REGEX = "REGEX"
	tok_MATCH = "MATCH"
	tok_END   = "END"

	tok_DEFINE = "="
	tok_DEFEND = ";"

	tok_OR    = "|"
	tok_MAYBE = "?"
	tok_MAYB1 = "+"
	tok_MATB0 = "*"

	tok_LPARENT = "("
	tok_RPARENT = ")"
)

type token struct {
	tokenType string
	literal   string
}

func newToken(tokenType string, literal string) token {
	return token{tokenType: tokenType, literal: literal}
}

type lexer struct {
	code string
	idx  int
}

func (mod *lexer) next() token {
	read := func() byte {
		mod.idx++
		return mod.code[mod.idx-1]
	}

	peek := func(count int) byte {
		return mod.code[mod.idx+count]
	}

	// skip white space characters.
	for mod.idx < len(mod.code) {
		if ch := peek(0); !u.IsSpace(ch) {
			break
		}

		read()
	}

	if mod.idx >= len(mod.code) {
		return newToken(tok_END, "END")
	}

	var tok token
	var ch byte = read()

	switch ch {
	case ';', '=', '|', '?', '+', '*', '(', ')':
		tok = newToken(string(ch), string(ch))

	case '\'':
		// match
		var match bytes.Buffer

		for mod.idx < len(mod.code) {
			if ch = read(); ch == '\\' {
				charMap := map[byte]byte{'\'': '\'', 'n': '\n', 'r': '\r', 't': '\t'}
				if ch, ok := charMap[read()]; ok == true {
					match.WriteByte(ch)
				} else {
					panic("\\?")
				}
			} else if ch != '\'' {
				match.WriteByte(ch)
			} else {
				break
			}
		}

		tok = newToken(tok_MATCH, match.String())

	case '$':
		// ident
		var ident bytes.Buffer

		for u.IsAlnum(peek(0)) || peek(0) == '_' {
			ident.WriteByte(read())
		}

		for mod.idx < len(mod.code) {
			if ch = peek(0); u.IsAlnum(ch) || ch == '_' {
				ident.WriteByte(read())
			} else {
				break
			}
		}

		tok = newToken(tok_IDENT, ident.String())

	default:
		// regex
		var regex bytes.Buffer

		if ch == 'R' && peek(0) == '<' {
			read()

			for mod.idx < len(mod.code) {
				if ch = read(); ch != '>' || peek(0) != 'R' {
					regex.WriteByte(ch)
				} else {
					break
				}
			}

			tok = newToken(tok_REGEX, regex.String())
		}
	}

	return tok
}

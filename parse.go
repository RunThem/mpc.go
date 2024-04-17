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
		return mod.code[mod.idx]
	}

	peek := func(count int) byte {
		return mod.code[mod.idx+count]
	}

	for mod.idx < len(mod.code) {
		ch := read()

		if u.IsSpace(ch) {
			continue
		}

		switch ch {
		case ';', '=', '|', '?', '+', '*', '(', ')':
			return newToken(string(ch), string(ch))

		case '\'':
			// match
			var match bytes.Buffer

			for {
				ch = read()

				if ch == '\\' {
					switch peek(0) {
					case '\'':
						match.WriteByte('\'')
					case 'n':
						match.WriteByte('\n')
					case 'r':
						match.WriteByte('\r')
					case 't':
						match.WriteByte('\t')
					default:
						panic("\\?")
					}

					read()
				} else if ch == '\'' {
					break
				} else {
					match.WriteByte(ch)
				}
			}

			return newToken(tok_MATCH, match.String())

		case '$':
			// ident
			var ident bytes.Buffer

			for u.IsAlnum(peek(0)) || peek(0) == '_' {
				ident.WriteByte(read())
			}

			return newToken(tok_IDENT, ident.String())

		default:
			// regex
			var regex bytes.Buffer

			if ch == 'R' && peek(0) == '<' {
				read()

				for {
					ch = read()
					if ch == '>' && peek(0) == 'R' {
						break
					}

					regex.WriteByte(ch)
				}

				return newToken(tok_REGEX, regex.String())
			}
		}
	}

	return newToken(tok_IDENT, mod.code[mod.idx:])
}

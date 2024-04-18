package mpc

import (
	"bytes"

	"github.com/RunThem/u"
)

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

/*
 * lexer
 */
type lexer struct {
	code string
	idx  int
}

func (mod *lexer) next() token {
	next := func() byte {
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

		next()
	}

	if mod.idx >= len(mod.code) {
		return newToken(tok_END, "END")
	}

	var tok token
	var ch byte = next()

	switch ch {
	case ';', '=', '|', '?', '+', '*', '(', ')':
		tok = newToken(string(ch), string(ch))

	case '\'':
		// match
		var match bytes.Buffer

		for mod.idx < len(mod.code) {
			if ch = next(); ch == '\\' {
				charMap := map[byte]byte{'\'': '\'', 'n': '\n', 'r': '\r', 't': '\t'}
				if ch, ok := charMap[next()]; ok == true {
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
			ident.WriteByte(next())
		}

		for mod.idx < len(mod.code) {
			if ch = peek(0); u.IsAlnum(ch) || ch == '_' {
				ident.WriteByte(next())
			} else {
				break
			}
		}

		tok = newToken(tok_IDENT, ident.String())

	default:
		// regex
		var regex bytes.Buffer

		if ch == 'R' && peek(0) == '<' {
			next()

			for mod.idx < len(mod.code) {
				if ch = next(); ch != '>' || peek(0) != 'R' {
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

func (mod *lexer) peek() token {
	idx := mod.idx
	tok := mod.next()
	mod.idx = idx

	return tok
}

/*
 * parser
 */

// func New(code string) *node {
// 	lex := lexer{code: code, idx: 0}
// 	defTable := make(map[string]*node)

// 	parse(&lex, defTable)
// }

// func parse(lex *lexer, defTable map[string]*node) {
// 	tok := lex.next()

// define *Mpc
// 	if tok.tokenType == tok_IDENT && lex.peek().tokenType == tok_DEFINE {
// 		defTable[tok.literal] = newMpc(m_and, tok.literal, "", nil)
// 	}
// }

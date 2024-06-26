package mpc

import (
	"testing"
)

func TestNext(t *testing.T) {
	code := `
    $factor = R<[0-9]*>R | '(' $expr ')'       ;
    $term   = $factor (('*' | '/') $factor)* ;
    $expr   = $term (('+' | '-') $term)      ;
`

	lex := &lexer{code: code, idx: 0}

	for tok := lex.next(); tok.tokenType != tok_END; tok = lex.next() {
		t.Log(tok)
	}
}

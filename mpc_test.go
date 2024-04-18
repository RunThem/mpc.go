package mpc

import (
	"testing"
)

func TestNew(t *testing.T) {
	/*
	 * $Factor = R<[0-9]*>R | '(' $Term ')'
	 * $Term   = <Factor> (('*' | '/') $Factor)*
	 */

	Term := newNode(m_and, "Term")
	Factor := newNode(m_or, "Factor")

	factor := newNode(m_and, "factor")
	factor.add(match("(", nil), Term, match(")", nil))

	Factor.add(regex("[0-9]*", nil), factor)

	term := newNode(m_maybe, "Term")

	op := newNode(m_or, "term_op")
	op.add(match("*", nil), match("/", nil))

	term.add(op, Factor)

	Term.add(Factor, term)

	Term.demp()
}

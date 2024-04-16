package mpc

import "testing"

func TestnewMpc(t *testing.T) {
	mpc := newMpc(m_match, "hello", "hello")

	_ = mpc
}

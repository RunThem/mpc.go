package mpc

import (
	"github.com/RunThem/u"
)

type mtag int

const (
	_ mtag = iota
	m_regex
	m_match

	m_and
	m_or
	m_maybe // +
	m_mayb1 //
	m_mayb0 // ?
)

type Mpc struct {
	tag   mtag
	name  string
	match string

	child *u.Vec[*Mpc]
}

func newMpc(tag mtag, name string, match string) *Mpc {
	return &Mpc{tag: tag, name: name, match: match}
}

func (mod *Mpc) addMpc(mpc ...*Mpc) {
	for _, it := range mpc {
		mod.child.PutBack(it)
	}
}

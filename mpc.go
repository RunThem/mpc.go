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
	m_maybe // ? 零次或一次
	m_mayb1 // + 一次或多次
	m_mayb0 // * 零次或多次
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

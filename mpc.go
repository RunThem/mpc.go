package mpc

import (
	"regexp"
	"strings"

	"github.com/RunThem/u"
)

type mtag int

const (
	_ mtag = iota
	m_ref

	m_regex
	m_match

	m_and
	m_or
	m_maybe // ? 零次或一次
	m_mayb1 // + 一次或多次
	m_mayb0 // * 零次或多次
)

type node struct {
	tag  mtag
	name string

	match  string
	regex  *regexp.Regexp
	parent *node
	*u.Vec[*node]
}

type Mpc struct {
	comment string
	root    *node
	refs    map[string]*node
	cur     *node

	input string
	idx   int
	col   int
	row   int
}

type Ast struct {
	content string

	*u.Vec[*Ast]
}

func New(comment string) *Mpc {
	return &Mpc{comment: comment, root: nil, refs: make(map[string]*node)}
}

func match(match string, parent *node) *node {
	return &node{tag: m_match, match: match, parent: parent}
}

func regex(regex string, parent *node) *node {
	return &node{tag: m_regex, match: regex, parent: parent}
}

func addChild(root *node, childs ...*node) {
	for _, it := range childs {
		root.PutBack(it)
	}
}

func (mod *Mpc) parseMatch(root *node) *Ast {
	if strings.HasPrefix(mod.input[mod.idx:], root.match) {
		mod.idx += len(root.match)
		return &Ast{content: root.match}
	}

	return nil
}

func (mod *Mpc) parseRegex(root *node) *Ast {
	result := root.regex.FindString(mod.input[mod.idx:])
	if result == "" {
		return nil
	}

	mod.idx += len(result)

	return &Ast{content: result}
}

func (mod *Mpc) parseAndExpr(root *node) *Ast {
	idx := mod.idx
	ast := &Ast{}

	for _, it := range root.Range(true) {
	}
}

func (mod *Mpc) parseOrExpr(root *node) *Ast {

}

func (mod *Mpc) Def(name string, other ...any) *Mpc {
	for _, it := range other {
		switch it.(type) {
		case string:

		}
	}

	return mod
}

/*

func (mod *Mpc) R(regex ...string) *Mpc {
	return mod
}

func (mod *Mpc) M(match ...string) *Mpc {
	return mod
}

func (mod *Mpc) Or() *Mpc {
	return mod
}

func (mod *Mpc) At(varname string) *Mpc {
	return mod
}

func (mod *Mpc) E() *Mpc {
	return mod
}

func (mod *Mpc) Maybe(mode byte) *Mpc {
	return mod
}

func (mod *Mpc) Sub(mpc *Mpc) *Mpc {
	return mod
}

func newMpc(tag mtag, name string, match string, parent *node) *node {
	return &node{tag: tag, name: name, match: match, parent: parent}
}

func (mod *node) addMpc(mpc ...*node) {
	for _, it := range mpc {
		mod.childs.PutBack(it)
	}
}

func main() {
	// $factor = R<[0-9]*>R | '(' $expr ')'       ;
	// $term   = $factor (('*' | '/') $factor)* ;
	// $expr   = $term (('+' | '-') $term)      ;

	calc := New("calc")
	calc.Def("factor").R("[0-9]*").Or().M("(").At("expr").M(")")
	calc.Def("term").At("factor").
		E().
		E().
		M("*", "/").
		E().At("factor").
		E().Maybe('*')

	calc.Def("term").At("factor").
		E().
		E().
		M("*", "/").
		E().At("factor").
		E().Maybe('*')

	calc.Def("term").At("factor").Sub(
		calc.Sub(
			calc.M("*", "/"),
		).At("factor"),
	).Maybe('*')

	// calc.Def("expr", m.R("[0-9]*"), m.Or, m.M("("), )
}

*/

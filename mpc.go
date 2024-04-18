package mpc

import (
	"fmt"
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

func (mod mtag) String() string {
	m := map[mtag]string{m_ref: "m_ref", m_regex: "m_regex", m_match: "m_match", m_and: "m_and", m_or: "m_or", m_maybe: "m_maybe", m_mayb1: "m_mayb1", m_mayb0: "m_mayb0"}

	return m[mod]
}

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

func newNode(tag mtag, name string) *node {
	return &node{
		tag:    tag,
		name:   name,
		match:  "",
		regex:  nil,
		parent: nil,
		Vec:    u.NewVec[*node](nil),
	}
}

func (mod *node) add(other ...*node) {
	for _, it := range other {
		mod.PutBack(it)
	}
}

func (mod *node) demp() {
	history := make(map[string]bool)

	var dump func(*node, int)
	dump = func(root *node, depth int) {
		fmt.Printf("%s", strings.Repeat("  ", depth))

		if len(root.name) != 0 {
			if _, ok := history[root.name]; ok {
				fmt.Printf("{%s:%p}\n", root.name, root)
				return
			}

			fmt.Printf("%s:%p ", root.name, root)
			history[root.name] = true
		}

		fmt.Printf("{%s '%s'", root.tag, root.match)

		if root.Vec == nil || root.Len() == 0 {
			fmt.Println("")
			return
		}

		fmt.Printf(" size %d\n", root.Len())
		for _, it := range root.Range(true) {
			dump(it, depth+1)
		}
	}

	dump(mod, 0)
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

func (mod *Mpc) parseAnd(root *node) *Ast {
	idx := mod.idx
	ast := &Ast{}

	for _, it := range root.Range(true) {
		if mod.idx == len(mod.input) {
			break
		}

		a := mod.parse(it)

		if a == nil {
			if it.tag == m_mayb0 || it.tag == m_mayb1 {
				continue
			}

			mod.idx = idx
			return nil
		}

		if it.tag == m_maybe || it.tag == m_mayb1 {
			for _, sast := range a.Range(true) {
				ast.PutBack(sast)
			}
		} else {
			ast.PutBack(a)
		}
	}

	return ast
}

func (mod *Mpc) parseOr(root *node) *Ast {
	idx := mod.idx

	for _, it := range root.Range(true) {
		if mod.idx == len(mod.input) {
			break
		}

		a := mod.parse(it)

		if a == nil {
			return a
		}

		mod.idx = idx
	}

	return nil
}

func (mod *Mpc) parseMaybe(root *node) *Ast {
	vec := u.NewVec[*Ast](nil)
	ast := &Ast{}

	for mod.idx != len(mod.input) {
		idx := mod.idx
		j := 0

		for i, it := range root.Range(true) {
			a := mod.parse(it)
			if a == nil {
				j = i
				break
			}

			vec.PutBack(a)
		}

		if j == root.Len() {
			for _, it := range vec.Range(true) {
				ast.PutBack(it)
			}
		} else {
			mod.idx = idx
			break
		}
	}

	return ast
}

func (mod *Mpc) parseMayb1(root *node) *Ast {
	vec := u.NewVec[*Ast](nil)
	ast := &Ast{}
	one := false

	for mod.idx != len(mod.input) {
		idx := mod.idx
		j := 0

		for i, it := range root.Range(true) {
			a := mod.parse(it)
			if a == nil {
				j = i
				break
			}

			vec.PutBack(a)
		}

		if j == root.Len() {
			one = true
			for _, it := range vec.Range(true) {
				ast.PutBack(it)
			}
		} else {
			mod.idx = idx
			break
		}
	}

	if !one {
		return nil
	}

	return ast
}

func (mod *Mpc) parseMayb0(root *node) *Ast {
	return nil
}

func (mod *Mpc) parse(root *node) *Ast {
	switch root.tag {
	case m_regex:
		return mod.parseRegex(root)
	case m_match:
		return mod.parseMatch(root)

	case m_and:
		return mod.parseAnd(root)
	case m_or:
		return mod.parseOr(root)

	case m_maybe:
		return mod.parseMaybe(root)
	case m_mayb1:
		return mod.parseMayb1(root)
	case m_mayb0:
		return mod.parseMayb0(root)
	}

	return nil
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

package compiler

import (
	"strings"

	"github.com/johnfrankmorgan/gazebo/debug"
	"github.com/johnfrankmorgan/gazebo/errors"
)

func parse(source string) *sexpr {
	tokens := tokenize(source)

	if debug.Enabled() {
		tokens.dump()
	}

	expr, _ := _parse(tokens)
	return expr
}

func _parse(tokens tokens) (*sexpr, int) {
	parser := parser{tokens: tokens}
	expr := parser.parse()
	return expr, parser.position
}

type sexpr struct {
	children []*sexpr
	token    token
}

func (m *sexpr) dump(depth int) {
	indent := strings.Repeat(" ", depth*2)

	if m.atom() || m.token.is(tkbracketopen, tkbracketclose) {
		debug.Printf("%s%s(%v)\n", indent, m.token.typ.name(), m.token.value)
		return
	}

	debug.Printf("%s(\n", indent)

	for _, expr := range m.children {
		expr.dump(depth + 1)
	}

	debug.Printf("%s)\n", indent)
}

func (m *sexpr) atom() bool {
	return m.token.atom()
}

type parser struct {
	tokens   tokens
	position int
}

func (m *parser) unexpectedeof() *sexpr {
	errors.ErrEOF.Panic("unexpected eof at token offset %d", m.position)
	return nil
}

func (m *parser) finished() bool {
	return m.position >= len(m.tokens)
}

func (m *parser) peek() token {
	return m.tokens[m.position]
}

func (m *parser) next() token {
	token := m.tokens[m.position]
	m.position++
	return token
}

func (m *parser) subexpr(start int, opener, closer tokentype) []token {
	errors.ErrParse.Expect(
		m.tokens[start].is(opener),
		"parse error: expected token type %s, got %s",
		opener.name(),
		m.tokens[start].typ.name(),
	)

	depth := 0

	for idx, token := range m.tokens[start:] {
		if token.is(opener) {
			depth++
		} else if token.is(closer) {
			depth--
			if depth == 0 {
				return m.tokens[start+1 : start+idx]
			}
		}
	}

	errors.ErrEOF.Panic("expecting %s near token offset %d", closer.name(), start)
	return nil
}

func (m *parser) parse() *sexpr {
	expr := new(sexpr)

	for !m.finished() {
		token := m.next()

		if token.is(tkwhitespace, tkcomment) {
			continue
		}

		if token.is(tkbracketopen) {
			subexpr, pos := _parse(m.subexpr(m.position-1, tkbracketopen, tkbracketclose))
			children := []*sexpr{&sexpr{token: token}}
			children = append(children, subexpr.children...)
			children = append(children, &sexpr{token: m.tokens[m.position+pos]})

			expr.children = append(expr.children, &sexpr{
				children: children,
			})

			m.position += pos + 1
			continue
		}

		if !token.is(tkparenopen) {
			expr.children = append(expr.children, &sexpr{
				token: token,
			})

			continue
		}

		subexpr, pos := _parse(m.subexpr(m.position-1, tkparenopen, tkparenclose))
		expr.children = append(expr.children, subexpr)

		m.position += pos + 1
	}

	return expr
}

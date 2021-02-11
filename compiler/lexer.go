package compiler

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/johnfrankmorgan/gazebo/assert"
)

func tokenize(source string) tokens {
	var tokens tokens

	lexer := lexer{source: []byte(strings.TrimSpace(source))}

	for !lexer.finished() {
		tokens = append(tokens, lexer.lex())
	}

	return tokens
}

type tokens []token

func (m tokens) dump() {
	for idx, tk := range m {
		fmt.Fprintf(
			os.Stderr,
			"%6d: %3d :: %16s :: %q\n",
			idx,
			tk.typ,
			tk.typ.name(),
			tk.value,
		)
	}
}

type tokentype int

const (
	tkinvalid tokentype = iota
	tkparenopen
	tkparenclose
	tkcomment
	tkwhitespace
	tkstring
	tkident
	tknumber
)

func (typ tokentype) valid() bool {
	return int(typ) > 0 && int(typ) <= int(tknumber)
}

func (typ tokentype) name() string {
	names := map[tokentype]string{
		tkinvalid:    "tkinvalid",
		tkparenopen:  "tkparenopen",
		tkparenclose: "tkparenclose",
		tkcomment:    "tkcomment",
		tkwhitespace: "tkwhitespace",
		tkstring:     "tkstring",
		tkident:      "tkident",
		tknumber:     "tknumber",
	}

	if name, ok := names[typ]; ok {
		return name
	}

	return "tkunknown"
}

type token struct {
	typ   tokentype
	value string
}

func (m *token) is(types ...tokentype) bool {
	for _, typ := range types {
		if m.typ == typ {
			return true
		}
	}

	return false
}

type lexer struct {
	source   []byte
	position int
	buffer   bytes.Buffer
}

func (m *lexer) unexpectedeof() token {
	assert.Unreached("unexpected eof at byte offset %d", m.position)
	return m.token(tkinvalid)
}

func (m *lexer) finished() bool {
	return m.position >= len(m.source)
}

func (m *lexer) peek() rune {
	ch, _ := utf8.DecodeRune(m.source[m.position:])
	return ch
}

func (m *lexer) next() rune {
	ch, width := utf8.DecodeRune(m.source[m.position:])
	m.buffer.WriteRune(ch)
	m.position += width
	return ch
}

func (m *lexer) token(typ tokentype) token {
	tk := token{typ: typ, value: m.buffer.String()}
	m.buffer.Reset()
	return tk
}

func (m *lexer) isdigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func (m *lexer) isalpha(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func (m *lexer) isidentchar(ch rune) bool {
	if m.isalpha(ch) || m.isdigit(ch) {
		return true
	}

	identchars := []rune{
		'!', '@', '£', '$', '%', '^', '&', '*',
		'-', '_', '+', '=', '<', '>', '?', '/',
		'.', '~', ':', ';',
	}

	for _, identch := range identchars {
		if identch == ch {
			return true
		}
	}

	return false
}

func (m *lexer) isnewline(ch rune) bool {
	return ch == '\n'
}

func (m *lexer) iswhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || m.isnewline(ch)
}

func (m *lexer) line(typ tokentype) token {
	for !m.finished() {
		if m.isnewline(m.peek()) {
			m.next()
			return m.token(typ)
		}

		m.next()
	}

	return m.unexpectedeof()
}

func (m *lexer) lstring() token {
	for !m.finished() {
		ch := m.peek()
		if ch == '"' {
			m.next()
			return m.token(tkstring)
		}

		m.next()
	}

	return m.unexpectedeof()
}

func (m *lexer) lnumber() token {
	var isfloat bool

	for !m.finished() {
		ch := m.peek()
		if ch == '.' && !isfloat {
			m.next()
			isfloat = true
			continue
		}

		if !m.isdigit(ch) {
			return m.token(tknumber)
		}

		m.next()
	}

	return m.unexpectedeof()
}

func (m *lexer) lident() token {
	for !m.finished() {
		ch := m.peek()
		if !m.isidentchar(ch) {
			return m.token(tkident)
		}

		m.next()
	}

	return m.unexpectedeof()
}

func (m *lexer) lwhitespace() token {
	for !m.finished() {
		ch := m.peek()
		if !m.iswhitespace(ch) {
			return m.token(tkwhitespace)
		}

		m.next()
	}

	return m.unexpectedeof()
}

func (m *lexer) lex() token {
	ch := m.next()

	switch ch {
	case '(':
		return m.token(tkparenopen)

	case ')':
		return m.token(tkparenclose)

	case ';':
		return m.line(tkcomment)

	case '"':
		return m.lstring()
	}

	if m.isdigit(ch) {
		return m.lnumber()
	}

	if m.isidentchar(ch) {
		return m.lident()
	}

	if m.iswhitespace(ch) {
		return m.lwhitespace()
	}

	assert.Unreached("unexpected rune %#U", ch)
	return m.token(tkinvalid)
}

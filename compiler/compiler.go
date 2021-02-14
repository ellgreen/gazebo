package compiler

import (
	"strconv"

	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/compiler/op"
	"github.com/johnfrankmorgan/gazebo/debug"
)

// Code is a slice containing executable bytecode
type Code []op.Instruction

// Dump prints a Code's formatted bytecode to stderr
func (m Code) Dump() {
	for idx, ins := range m {
		debug.Printf(
			"%6d %18s (0x%02x) %v\n",
			idx,
			ins.Opcode.Name(),
			int(ins.Opcode),
			ins.Arg,
		)
	}
}

// Compile compiles gazebo code into a Code object
func Compile(source string) Code {
	var (
		code     Code
		compiler compiler
	)

	expr := parse(source)

	if debug.Enabled() {
		expr.dump(0)
	}

	assert.False(expr.atom())

	for _, expr := range expr.children {
		code = append(code, compiler.compile(expr)...)
	}

	if debug.Enabled() {
		code.Dump()
	}

	return code
}

type compiler struct{}

func (m *compiler) compile(expr *sexpr) Code {
	if expr.atom() {
		return m.atom(expr.token)
	}

	if len(expr.children) == 0 {
		return Code{}
	}

	switch expr.children[0].token.value {
	case "let":
		assert.Len(expr.children, 3)
		assert.True(expr.children[1].atom())
		return append(
			m.compile(expr.children[2]),
			op.StoreName.Ins(expr.children[1].token.value),
		)

	case "if":
		assert.Len(expr.children, 4)
		truepath := m.compile(expr.children[2])
		falsepath := append(
			m.compile(expr.children[3]),
			op.RelJump.Ins(len(truepath)),
		)
		code := append(
			m.compile(expr.children[1]),
			op.RelJumpIfTrue.Ins(len(falsepath)),
		)
		code = append(code, falsepath...)
		return append(code, truepath...)

	case "while":
		assert.Len(expr.children, 3)
		body := m.compile(expr.children[2])
		cond := append(
			m.compile(expr.children[1]),
			op.RelJumpIfFalse.Ins(len(body)+1),
		)
		body = append(
			body,
			op.RelJump.Ins(-len(body)-len(cond)-1),
		)
		return append(cond, body...)

	case "fun":
		if len(expr.children) == 3 {
			params := []string{}
			for _, param := range expr.children[1].children {
				assert.True(
					param.atom() && param.token.is(tkident),
					"function parameters must be identifiers",
				)

				params = append(params, param.token.value)
			}
			code := Code{
				op.PushValue.Ins(params),
				op.PushValue.Ins(m.compile(expr.children[2])),
			}
			return append(code, op.MakeFunc.Ins(len(params)))

		} else if len(expr.children) == 4 {
			assert.True(
				expr.children[1].atom() && expr.children[1].token.is(tkident),
				"function name must be a valid identifier",
			)
			name := expr.children[1].token.value
			params := expr.children[2]
			body := expr.children[3]
			return append(
				m.compile(&sexpr{children: []*sexpr{{token: token{typ: tkident, value: "fun"}}, params, body}}),
				op.StoreName.Ins(name),
			)
		}

		assert.Unreached("fun keyword should contain 3 or 4 children, got %d", len(expr.children))

	case "load":
		assert.True(len(expr.children) >= 2)
		code := Code{}

		for _, expr := range expr.children[1:] {
			assert.True(
				expr.atom() && expr.token.is(tkident),
				"load parameters must be identifiers",
			)

			code = append(code, op.LoadModule.Ins(expr.token.value))
		}

		return code
	}

	if !expr.children[0].atom() {
		code := Code{}

		for _, expr := range expr.children {
			code = append(code, m.compile(expr)...)
		}

		return code
	}

	function := m.compile(expr.children[0])
	argc := 0

	for _, arg := range expr.children[1:] {
		function = append(function, m.compile(arg)...)
		argc++
	}

	return append(function, op.CallFunc.Ins(argc))
}

func (m *compiler) atom(tk token) Code {
	switch tk.typ {
	case tknumber:
		value, err := strconv.ParseFloat(tk.value, 64)
		assert.Nil(err)

		return Code{op.LoadConst.Ins(value)}

	case tkstring:
		return Code{op.LoadConst.Ins(tk.value[1 : len(tk.value)-1])}

	case tkident:
		return Code{op.LoadName.Ins(tk.value)}
	}

	assert.Unreached("invalid atom: %#v", tk)
	return nil
}

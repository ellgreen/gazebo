package compiler

import (
	"strconv"

	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/compiler/op"
	"github.com/johnfrankmorgan/gazebo/debug"
	"github.com/johnfrankmorgan/gazebo/errors"
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
func Compile(source string) (code Code, err error) {
	defer func() {
		recovered := recover()

		if recovered == nil {
			return
		}

		if gerr, ok := recovered.(*errors.Error); ok {
			err = gerr
			return
		}

		panic(recovered)
	}()

	var compiler compiler

	expr := parse(source)

	if debug.Enabled() {
		expr.dump(0)
	}

	errors.ErrParse.Expect(!expr.atom(), "unexpected atom, expecting list")

	for _, expr := range expr.children {
		code = append(code, compiler.compile(expr)...)
	}

	if debug.Enabled() {
		code.Dump()
	}

	return
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
		errors.ErrParse.ExpectLen(expr.children, 3, "invalid let syntax")

		errors.ErrParse.Expect(
			expr.children[1].token.is(tkident),
			"expected %s, got %s",
			tkident.name(),
			expr.children[1].token.typ.name(),
		)

		return append(
			m.compile(expr.children[2]),
			op.StoreName.Ins(expr.children[1].token.value),
		)

	case "if":
		errors.ErrParse.ExpectLen(expr.children, 4, "invalid if syntax")
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
		errors.ErrParse.ExpectLen(expr.children, 3, "invalid while syntax")
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
				errors.ErrParse.Expect(
					param.token.is(tkident),
					"function parameters must be identifiers, got %s",
					param.token.typ.name(),
				)

				params = append(params, param.token.value)
			}
			code := Code{
				op.PushValue.Ins(params),
				op.PushValue.Ins(m.compile(expr.children[2])),
			}
			return append(code, op.MakeFunc.Ins(len(params)))

		} else if len(expr.children) == 4 {
			errors.ErrParse.Expect(
				expr.children[1].token.is(tkident),
				"function name must be a valid identifier, got %s",
				expr.children[1].token.typ.name(),
			)
			name := expr.children[1].token.value
			params := expr.children[2]
			body := expr.children[3]
			return append(
				m.compile(&sexpr{children: []*sexpr{{token: token{typ: tkident, value: "fun"}}, params, body}}),
				op.StoreName.Ins(name),
			)
		}

		errors.ErrParse.Panic("invalid fun syntax")

	case "load":
		errors.ErrParse.ExpectAtLeast(expr.children, 2, "invalid load syntax")
		code := Code{}

		for _, expr := range expr.children[1:] {
			errors.ErrParse.Expect(
				expr.token.is(tkident),
				"load parameters must be identifiers, got %s",
				expr.token.typ.name(),
			)

			code = append(code, op.LoadModule.Ins(expr.token.value))
		}

		return code

	case "list":
		errors.ErrParse.ExpectLen(expr.children, 2, "invalid list syntax")

		code := Code{}
		length := 0

		for _, expr := range expr.children[1].children {
			code = append(code, m.compile(expr)...)
			length++
		}

		return append(code, op.MakeList.Ins(length))
	}

	if expr.children[0].token.is(tkbracketopen) {
		errors.ErrParse.ExpectLen(expr.children, 3, "invalid index syntax")
		errors.ErrParse.Expect(
			expr.children[2].token.is(tkbracketclose),
			"expected %s, got %s",
			tkbracketclose.name(),
			expr.children[2].token.typ.name(),
		)
		if expr.children[1].token.is(tkident) && expr.children[1].token.value[0] == '.' {
			return Code{op.AttributeGet.Ins(expr.children[1].token.value[1:])}
		}
		return append(m.compile(expr.children[1]), op.IndexGet.Ins(nil))
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
		if len(arg.children) > 0 && arg.children[0].token.is(tkbracketopen) {
			continue
		}
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

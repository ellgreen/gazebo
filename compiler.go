package gazebo

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/johnfrankmorgan/gazebo/assert"
)

// Opcode is an opcode recognised by the gazebo VM
type Opcode int

// Enumeration of available Opcodes
const (
	_ Opcode = iota
	OpLoadConst
	OpStoreName
	OpLoadName
	OpCallFunc
	OpRelJump
	OpRelJumpIfTrue
	OpRelJumpIfFalse
)

// Ins creates an Instruction for an Opcode
func (op Opcode) Ins(arg interface{}) Instruction {
	return Instruction{Opcode: op, Arg: arg}
}

// Name returns an Opcode's name
func (op Opcode) Name() string {
	names := map[Opcode]string{
		OpLoadConst:      "OpLoadConst",
		OpStoreName:      "OpStoreName",
		OpLoadName:       "OpLoadName",
		OpCallFunc:       "OpCallFunc",
		OpRelJump:        "OpRelJump",
		OpRelJumpIfTrue:  "OpRelJumpIfTrue",
		OpRelJumpIfFalse: "OpRelJumpIfFalse",
	}

	if name, ok := names[op]; ok {
		return name
	}

	return "OpUnknown"
}

// Instruction is a struct containing an Opcode
// and an optional argument
type Instruction struct {
	Opcode Opcode
	Arg    interface{}
}

// Code is a slice containing executable bytecode
type Code []Instruction

// Dump prints a Code's formatted bytecode to stderr
func (m Code) Dump() {
	for idx, ins := range m {
		fmt.Fprintf(
			os.Stderr,
			"%6d %18s (%d) %v\n",
			idx,
			ins.Opcode.Name(),
			ins.Opcode,
			ins.Arg,
		)
	}
}

// Compile compiles gazebo code into a Code object
func Compile(source string) Code {
	var compiler Compiler

	return compiler.Compile(source)
}

// Compiler is a struct to compile a sequence of
// gazebo tokens
type Compiler struct{}

// Compile compiles source code into a Code object
func (m *Compiler) Compile(source string) Code {
	var (
		code   Code
		parser parser
	)

	expr := parser.parse(parser.split(source))

	for _, expr := range expr.children {
		code = append(code, m.compile(expr)...)
	}

	return code
}

func (m *Compiler) compile(expr *sexpr) Code {
	if expr.isAtom() {
		return m.atom(expr.value)
	}

	if len(expr.children) == 0 {
		return Code{}
	}

	switch expr.children[0].value {
	case "let":
		assert.Len(expr.children, 3)
		assert.True(expr.children[1].isAtom())
		return append(
			m.compile(expr.children[2]),
			OpStoreName.Ins(expr.children[1].value),
		)

	case "if":
		assert.Len(expr.children, 4)
		truepath := m.compile(expr.children[2])
		falsepath := append(
			m.compile(expr.children[3]),
			OpRelJump.Ins(len(truepath)),
		)
		code := append(
			m.compile(expr.children[1]),
			OpRelJumpIfTrue.Ins(len(falsepath)),
		)
		code = append(code, falsepath...)
		return append(code, truepath...)

	case "while":
		assert.Len(expr.children, 3)
		body := m.compile(expr.children[2])
		cond := append(
			m.compile(expr.children[1]),
			OpRelJumpIfFalse.Ins(len(body)+1),
		)
		body = append(
			body,
			OpRelJump.Ins(-len(body)-len(cond)-1),
		)
		return append(cond, body...)
	}

	if !expr.children[0].isAtom() {
		code := make(Code, 0)

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

	return append(function, OpCallFunc.Ins(argc))
}

func (m *Compiler) atom(value string) Code {
	var regexes = struct {
		numbers *regexp.Regexp
		strings *regexp.Regexp
		idents  *regexp.Regexp
	}{
		numbers: regexp.MustCompile(`^-?[0-9]+(.[0-9]+)?$`),
		strings: regexp.MustCompile(`^".*"$`),
		idents:  regexp.MustCompile(`^[a-zA-Z!@$%^&*\/?<>_=+~-]+$`),
	}

	switch true {
	case regexes.numbers.MatchString(value):
		value, err := strconv.ParseFloat(value, 64)
		assert.Nil(err)

		return Code{OpLoadConst.Ins(value)}

	case regexes.strings.MatchString(value):
		value = value[1 : len(value)-1]
		value = strings.ReplaceAll(value, "\\n", "\n")
		return Code{OpLoadConst.Ins(value)}

	case regexes.idents.MatchString(value):
		return Code{OpLoadName.Ins(value)}
	}

	assert.Unreached("invalid atom: %q", value)

	return nil
}

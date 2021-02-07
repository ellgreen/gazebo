package op

import "github.com/johnfrankmorgan/gazebo/gvalue"

// Opcode represents an opcode in the gazebo interpreter
type Opcode int

// Enumeration of Opcodes
const (
	Invalid Opcode = iota
	LoadConst
	StoreName
	LoadName
	CallFunc
)

func (op Opcode) String() string {
	names := map[Opcode]string{
		Invalid:   "Invalid",
		LoadConst: "LoadConst",
		StoreName: "StoreName",
		LoadName:  "LoadName",
		CallFunc:  "CallFunc",
	}

	return names[op]
}

// Instruction creates a new Instruction for the provided opcode
func (op Opcode) Instruction(arg gvalue.Instance) Instruction {
	return Instruction{Opcode: op, Arg: arg}
}

// Instruction is a struct containing an opcode and an optional argument
type Instruction struct {
	Opcode Opcode
	Arg    gvalue.Instance
}

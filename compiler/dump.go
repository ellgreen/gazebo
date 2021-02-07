package compiler

import (
	"github.com/johnfrankmorgan/gazebo/op"
	"github.com/kr/pretty"
)

func Dump(code []op.Instruction) {
	for idx, ins := range code {
		pretty.Printf("%6d %10s (%d) %# v\n",
			idx,
			ins.Opcode.String(),
			ins.Opcode,
			ins.Arg,
		)
	}
}

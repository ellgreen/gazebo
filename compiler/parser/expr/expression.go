package expr

import "github.com/johnfrankmorgan/gazebo/compiler/code"

type Expression interface {
	code.Compiler
}

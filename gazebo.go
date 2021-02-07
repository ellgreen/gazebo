package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/compiler"
	"github.com/johnfrankmorgan/gazebo/interpreter"
)

func main() {
	assert.True(len(os.Args) == 2)

	infile := os.Args[1]
	source, err := ioutil.ReadFile(infile)
	assert.Nil(err, "ioutil.ReadFile: %v", err)

	code := compiler.Compile(string(source))
	compiler.Dump(code)
	fmt.Print("\n\n")

	interpreter.Create()
	interpreter.The().Eval(code)
}

package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/compiler"
	"github.com/johnfrankmorgan/gazebo/interpreter"
)

var (
	file string
)

func main() {
	flag.StringVar(&file, "file", "", "file to run")
	flag.Parse()

	source, err := ioutil.ReadFile(file)
	assert.Nil(err, "ioutil.ReadFile: %v", err)

	code := compiler.Compile(string(source))
	compiler.Dump(code)
	fmt.Print("\n\n")

	interpreter.Create()
	interpreter.The().Eval(code)
}

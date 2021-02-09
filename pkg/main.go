package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/johnfrankmorgan/gazebo"
	"github.com/johnfrankmorgan/gazebo/assert"
)

func main() {
	assert.Len(os.Args, 2)

	source, err := ioutil.ReadFile(os.Args[1])
	assert.Nil(err)

	code := gazebo.Compile(string(source))
	code.Dump()

	fmt.Print("\n\n")

	gazebo.NewVM().Run(code)
}

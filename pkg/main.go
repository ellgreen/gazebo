package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/johnfrankmorgan/gazebo"
	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/compiler"
	"github.com/johnfrankmorgan/gazebo/debug"
)

func main() {
	debugging := flag.Bool("d", false, "enable debugging")

	flag.Parse()

	if *debugging {
		debug.Enable()
	}

	assert.Len(flag.Args(), 1)

	source, err := ioutil.ReadFile(flag.Args()[0])
	assert.Nil(err)

	code := compiler.Compile(string(source))

	if debug.Enabled() {
		fmt.Print("\n\n")
	}

	gazebo.NewVM().Run(code)
}

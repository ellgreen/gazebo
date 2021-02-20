package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/compiler"
	"github.com/johnfrankmorgan/gazebo/debug"
	"github.com/johnfrankmorgan/gazebo/errors"
	"github.com/johnfrankmorgan/gazebo/g"
	"github.com/johnfrankmorgan/gazebo/vm"
)

func main() {
	debugging := flag.Bool("d", false, "enable debugging")

	flag.Parse()

	if *debugging {
		debug.Enable()
	}

	if len(flag.Args()) == 0 {
		repl()
		return
	}

	source, err := ioutil.ReadFile(flag.Args()[0])
	assert.Nil(err)

	code, err := compiler.Compile(string(source))
	assert.Nil(err)

	if debug.Enabled() {
		debug.Printf("\n\n")
	}

	_, err = vm.New(flag.Args()[1:]...).Run(code)
	assert.Nil(err)
}

func repl() {
	var (
		buffer strings.Builder
		more   bool
	)

	vm := vm.New()

	for {
		if more {
			fmt.Printf("... ")
		} else {
			fmt.Printf(">>> ")
		}

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			buffer.WriteString(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "%s", err.Error())
			buffer.Reset()
			more = false
			continue
		}

		source := strings.TrimSpace(buffer.String())
		if len(source) > 0 && source[0] != '(' {
			source = fmt.Sprintf("(%s)", source)
		}

		code, err := compiler.Compile(source)
		if err != nil {
			if err == errors.ErrEOF {
				more = true
				buffer.WriteByte(' ')
			} else {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				buffer.Reset()
				more = false
			}

			continue
		}

		result, err := vm.Run(code)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}

		if result != nil && result.Type() != g.TypeNil {
			fmt.Printf("%v\n", result.Call(g.Protocols.Inspect, nil).Value())
		}

		buffer.Reset()
		more = false
	}
}

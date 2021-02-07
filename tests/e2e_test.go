package tests

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/johnfrankmorgan/gazebo/compiler"
	"github.com/johnfrankmorgan/gazebo/interpreter"
)

const TestScripts = "../tests/gaz"

func TestGazScripts(t *testing.T) {
	scripts, err := ioutil.ReadDir(TestScripts)
	if err != nil {
		t.Error(err)
	}

	for _, script := range scripts {
		t.Run(script.Name(), func(t *testing.T) {
			source, err := ioutil.ReadFile(path.Join(TestScripts, script.Name()))
			if err != nil {
				t.Error(err)
			}

			code := compiler.Compile(string(source))

			interpreter.Create()
			interpreter.The().Eval(code)
		})
	}
}

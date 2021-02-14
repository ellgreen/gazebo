package vm

import (
	"testing"

	"github.com/johnfrankmorgan/gazebo/compiler"
)

var (
	code compiler.Code
	vm   *VM
)

func init() {
	source := `
		(let N 20)
		(fun fib (n) (
			if (< n 2) (
				(+ n 0)
			) (
				(+ (fib (- n 2)) (fib (- n 1)))
			)
		))
		(fib N)
	`

	code = compiler.Compile(source)
	vm = New()
}

func BenchmarkVMRun(b *testing.B) {
	for n := 0; n < b.N; n++ {
		vm.Run(code)
	}
}

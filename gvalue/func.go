package gvalue

import "fmt"

type Func interface {
	Instance
	Name() string
	Call([]Instance) Instance
}

type builtin struct {
	name string
	f    func([]Instance) Instance
}

func (m builtin) ToString() string {
	return fmt.Sprintf("Builtin function %p", m.f)
}

func (m builtin) Interface() interface{} {
	return m.f
}

func (m builtin) Name() string {
	return m.name
}

func (m builtin) Call(args []Instance) Instance {
	return m.f(args)
}

func Builtin(name string, f func([]Instance) Instance) Func {
	return builtin{name: name, f: f}
}

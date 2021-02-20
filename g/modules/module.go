package modules

import (
	"fmt"

	"github.com/johnfrankmorgan/gazebo/g"
)

// All returns all available modules
func All() map[string]*Module {
	return map[string]*Module{
		"http": HTTP,
		"str":  Str,
		"time": Time,
	}
}

// Module defines g.Object values
type Module struct {
	Name   string
	Init   func()
	Values map[string]g.Object
}

// Load loads a Module's values into a g.Attributes
func (m *Module) Load(attrs *g.Attributes) {
	if m.Init != nil {
		m.Init()
	}

	for name, value := range m.Values {
		attrs.Set(fmt.Sprintf("%s.%s", m.Name, name), value)
	}
}

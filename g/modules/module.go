package modules

import (
	"fmt"

	"github.com/johnfrankmorgan/gazebo/g"
)

// Module defines g.Object values
type Module struct {
	Name   string
	Values map[string]g.Object
}

// Load loads a Module's values into a g.Attributes
func (m *Module) Load(attrs *g.Attributes) {
	for name, value := range m.Values {
		attrs.Set(fmt.Sprintf("%s.%s", m.Name, name), value)
	}
}

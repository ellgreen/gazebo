package gvalue

import "fmt"

// Bool is a gazebo boolean
type Bool struct {
	Value bool
}

func (m *Bool) ToString() string {
	return fmt.Sprintf("%v", m.Value)
}

func (m *Bool) Interface() interface{} {
	return m.Value
}

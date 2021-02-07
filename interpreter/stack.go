package interpreter

import (
	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/gvalue"
	"github.com/kr/pretty"
)

type Stack struct {
	values []gvalue.Instance
}

func (m *Stack) Push(value gvalue.Instance) {
	m.values = append(m.values, value)
}

func (m *Stack) Pop() gvalue.Instance {
	size := m.Size()

	if size > 0 {
		value := m.values[size-1]
		m.values = m.values[0 : size-1]
		return value
	}

	assert.Unreached("stack empty")
	return nil
}

func (m *Stack) Size() int {
	return len(m.values)
}

func (m *Stack) Dump() {
	for i := 0; i < m.Size(); i++ {
		pretty.Printf("%6d %# v\n", i, m.values[i])
	}
}

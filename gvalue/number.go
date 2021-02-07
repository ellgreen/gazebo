package gvalue

import (
	"fmt"
	"math"
)

// Number is the representation of integers / floats in gazebo
type Number struct {
	Value float64
}

func (m *Number) ToString() string {
	return fmt.Sprintf("%v", m.Value)
}

func (m *Number) ToBool() bool {
	return m.Value != 0
}

func (m *Number) Interface() interface{} {
	// Return int value if m.Value is a whole number
	if math.Mod(m.Value, 1) == 0 {
		return int(m.Value)
	}

	return m.Value
}

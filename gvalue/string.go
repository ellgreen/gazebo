package gvalue

// String is a gazebo string
type String struct {
	Value string
}

func (m *String) ToString() string {
	return m.Value
}

func (m *String) ToBool() bool {
	return len(m.Value) > 0
}

func (m *String) Interface() interface{} {
	return m.Value
}

package gvalue

// String is a gazebo string
type String struct {
	Value string
}

func (m *String) ToString() string {
	return m.Value
}

func (m *String) Interface() interface{} {
	return m.Value
}

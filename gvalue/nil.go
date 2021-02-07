package gvalue

// Nil is gazebo's null value
type Nil struct{}

func (m *Nil) ToString() string {
	return "nil"
}

func (m *Nil) Interface() interface{} {
	return nil
}

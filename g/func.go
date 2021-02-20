package g

import (
	"github.com/johnfrankmorgan/gazebo/errors"
)

// Func is the type of gazebo's builtin functions
type Func func(args Args) Object

// Args is used to pass arguments to builtin gazebo functions
type Args []Object

// Expects asserts that N arguments were received
func (m Args) Expects(n int) {
	errors.ErrRuntime.ExpectLen(m, n,
		"expected %d arguments, got %d", n, len(m))
}

// ExpectsAtLeast asserts that at least N arguments were received
func (m Args) ExpectsAtLeast(n int) {
	errors.ErrRuntime.ExpectAtLeast(m, n,
		"expected at least %d arguments, got %d", n, len(m))
}

// Self is a helper method to return the first argument
func (m Args) Self() Object {
	m.ExpectsAtLeast(1)
	return m[0]
}

// SelfWithArgs returns the first argument as an Object
// and the remaining args in an Args slice
func (m Args) SelfWithArgs() (Object, Args) {
	return m.Self(), m[1:]
}

// Values returns an array of interface{} values from an
// Args slice
func (m Args) Values() []interface{} {
	var ifaces []interface{}

	for _, arg := range m {
		ifaces = append(ifaces, arg.Value())
	}

	return ifaces
}

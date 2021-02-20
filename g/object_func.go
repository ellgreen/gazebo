package g

import "github.com/johnfrankmorgan/gazebo/compiler"

// ObjectFunc is the underlying type of user-defined functions in gazebo
type ObjectFunc struct {
	PartialObject
	params []string
	code   compiler.Code
	env    interface{}
}

// NewObjectFunc creates a new user-defined function
func NewObjectFunc(params []string, code compiler.Code, env interface{}) Object {
	return &ObjectFunc{
		PartialObject: PartialObject{typ: TypeFunc},
		params:        params,
		code:          code,
		env:           env,
	}
}

// Value satisfies the Object interface
func (m *ObjectFunc) Value() interface{} {
	return nil
}

// Call satisfies the Object interface
func (m *ObjectFunc) Call(method string, args Args) Object {
	return m.call(m, method, args)
}

// Params returns the function's parameter names
func (m *ObjectFunc) Params() []string {
	return m.params
}

// Code returns the function's bytecode
func (m *ObjectFunc) Code() compiler.Code {
	return m.code
}

// Env returns the function's environment
func (m *ObjectFunc) Env() interface{} {
	return m.env
}

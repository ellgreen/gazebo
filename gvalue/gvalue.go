package gvalue

import "github.com/johnfrankmorgan/gazebo/assert"

// Instance is the interface satisfied by all gazebo values
type Instance interface {
	ToString() string
	Interface() interface{}
}

// New converts a go value into the appropriate gazebo value
func New(value interface{}) Instance {
	if value == nil {
		return &Nil{}
	}

	switch value := value.(type) {
	case bool:
		return &Bool{value}

	case string:
		return &String{value}

	case int:
		return &Number{float64(value)}

	case float64:
		return &Number{value}
	}

	assert.Unreached("failed to convert value to gvalue.Instance: %# v", value)
	return nil
}

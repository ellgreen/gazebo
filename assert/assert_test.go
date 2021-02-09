package assert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrue(t *testing.T) {
	assert := assert.New(t)

	assert.PanicsWithError("Assertion failed: test True", func() {
		True(false, "test %s", "True")
	})
}

func TestFalse(t *testing.T) {
	assert := assert.New(t)

	assert.PanicsWithError("Assertion failed: test False", func() {
		False(true, "test False")
	})
}

func TestNil(t *testing.T) {
	assert := assert.New(t)

	assert.Panics(func() {
		Nil(1)
	})
}

func TestNotNil(t *testing.T) {
	assert := assert.New(t)

	assert.Panics(func() {
		NotNil(nil)
	})
}

func TestLen(t *testing.T) {
	assert := assert.New(t)

	assert.Panics(func() {
		Len([]string{""}, 2)
	})
}

func TestUnreached(t *testing.T) {
	assert := assert.New(t)

	assert.Panics(func() {
		Unreached()
	})
}

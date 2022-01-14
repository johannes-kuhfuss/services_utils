package misc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SliceContainsString_NotInSlice_Returns_False(t *testing.T) {
	source := []string{"New York", "Rio", "Tokyo"}

	in := SliceContainsString(source, "Paris")

	assert.EqualValues(t, false, in)
}

func Test_SliceContainsString_InSlice_Returns_True(t *testing.T) {
	source := []string{"New York", "Rio", "Tokyo"}

	in := SliceContainsString(source, "Rio")

	assert.EqualValues(t, true, in)
}

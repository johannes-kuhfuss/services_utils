package misc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceContainsStringNotInSliceReturnsFalse(t *testing.T) {
	source := []string{"New York", "Rio", "Tokyo"}

	in := SliceContainsString(source, "Paris")

	assert.EqualValues(t, false, in)
}

func TestSliceContainsStringInSliceReturnsTrue(t *testing.T) {
	source := []string{"New York", "Rio", "Tokyo"}

	in := SliceContainsString(source, "Rio")

	assert.EqualValues(t, true, in)
}

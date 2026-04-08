package misc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	source = []string{"New York", "Rio", "Tokyo"}
)

func TestSliceContainsStringNotInSliceReturnsFalse(t *testing.T) {
	in := SliceContainsString(source, "Paris")

	assert.EqualValues(t, false, in)
}

func TestSliceContainsStringInSliceReturnsTrue(t *testing.T) {
	in := SliceContainsString(source, "Rio")

	assert.EqualValues(t, true, in)
}

func TestSliceContainsStringCINotInSliceReturnsFalse(t *testing.T) {
	in := SliceContainsStringCI(source, "Paris")

	assert.EqualValues(t, false, in)
}

func TestSliceContainsStringCIInSliceLowerCaseReturnsTrue(t *testing.T) {
	in := SliceContainsStringCI(source, "rio")

	assert.EqualValues(t, true, in)
}

func TestSliceContainsStringCIInSliceUpperCaseReturnsTrue(t *testing.T) {
	in := SliceContainsStringCI(source, "RIO")

	assert.EqualValues(t, true, in)
}

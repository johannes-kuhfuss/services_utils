package date

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConst(t *testing.T) {
	assert.EqualValues(t, time.RFC3339, ApiDateLayout)
}

func Test_GetNowUtcString_ReturnsNoErr(t *testing.T) {
	now := GetNowUtcString()
	date, err := time.Parse(ApiDateLayout, now)
	assert.NotNil(t, date)
	assert.Nil(t, err)
	assert.EqualValues(t, date.Format(ApiDateLayout), now)
}

func Test_GetNowLocalString_WrongLocation_ReturnsErr(t *testing.T) {
	nowLocal, err := GetNowLocalString("wrong location")
	assert.Nil(t, nowLocal)
	assert.NotNil(t, err)
}

func Test_GetNowLocalString_EmptyLocation_ReturnsNoErr(t *testing.T) {
	nowLocal, err := GetNowLocalString("")
	assert.NotNil(t, nowLocal)
	assert.Nil(t, err)
	assert.True(t, IsValidTime(*nowLocal))
}

func Test_GetNowLocalString_WithLocation_ReturnsNoErr(t *testing.T) {
	nowLocal, err := GetNowLocalString("Europe/Berlin")
	assert.NotNil(t, nowLocal)
	assert.Nil(t, err)
	assert.True(t, IsValidTime(*nowLocal))
}

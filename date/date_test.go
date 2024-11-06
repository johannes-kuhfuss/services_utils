package date

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConst(t *testing.T) {
	assert.EqualValues(t, time.RFC3339, ApiDateLayout)
}

func TestGetNowUtcStringReturnsNoErr(t *testing.T) {
	now := GetNowUtcString()
	date, err := time.Parse(ApiDateLayout, now)
	assert.NotNil(t, date)
	assert.Nil(t, err)
	assert.EqualValues(t, date.Format(ApiDateLayout), now)
}

func TestGetNowLocalStringWrongLocationReturnsErr(t *testing.T) {
	nowLocal, err := GetNowLocalString("wrong location")
	assert.Nil(t, nowLocal)
	assert.NotNil(t, err)
}

func TestGetNowLocalStringEmptyLocationReturnsNoErr(t *testing.T) {
	nowLocal, err := GetNowLocalString("")
	assert.NotNil(t, nowLocal)
	assert.Nil(t, err)
	assert.True(t, IsValidTime(*nowLocal))
}

func TestGetNowLocalStringWithLocationReturnsNoErr(t *testing.T) {
	nowLocal, err := GetNowLocalString("Europe/Berlin")
	assert.NotNil(t, nowLocal)
	assert.Nil(t, err)
	assert.True(t, IsValidTime(*nowLocal))
}

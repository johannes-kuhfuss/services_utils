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

func TestIsValidTime(t *testing.T) {
	tests := []struct {
		name    string
		timeStr string
		want    bool
	}{
		{
			name:    "valid utc",
			timeStr: "2026-05-18T08:00:00Z",
			want:    true,
		},
		{
			name:    "valid offset",
			timeStr: "2026-05-18T10:00:00+02:00",
			want:    true,
		},
		{
			name:    "valid fractional seconds",
			timeStr: "2026-05-18T08:00:00.123Z",
			want:    true,
		},
		{
			name:    "date only",
			timeStr: "2026-05-18",
			want:    false,
		},
		{
			name:    "whitespace",
			timeStr: " 2026-05-18T08:00:00Z ",
			want:    false,
		},
		{
			name:    "invalid",
			timeStr: "not a date",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualValues(t, tt.want, IsValidTime(tt.timeStr))
		})
	}
}

func TestGetNowLocalStringEmptyLocationReturnsNoErr(t *testing.T) {
	nowLocal, err := GetNowLocalString("")
	assert.NotNil(t, nowLocal)
	assert.Nil(t, err)
	assert.True(t, IsValidTime(*nowLocal))
}

func TestGetNowLocalStringTrimsLocation(t *testing.T) {
	nowLocal, err := GetNowLocalString(" Europe/Berlin ")
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
